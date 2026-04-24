// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package apihttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/messaging"
	"sstpa-tool/backend/internal/metadata"
)

type messageResponse struct {
	MessageID                string         `json:"messageId"`
	Subject                  string         `json:"subject"`
	Body                     string         `json:"body,omitempty"`
	SentAt                   string         `json:"sentAt"`
	Sender                   string         `json:"sender"`
	SenderEmail              string         `json:"senderEmail"`
	MessageType              string         `json:"messageType"`
	IsRead                   bool           `json:"isRead"`
	RelatedNodeHIDs          []string       `json:"relatedNodeHids"`
	RelatedRelationshipTypes []string       `json:"relatedRelationshipTypes"`
	Properties               map[string]any `json:"properties"`
}

type unreadCountResponse struct {
	UnreadCount int64 `json:"unreadCount"`
}

type directMessageRequest struct {
	Actor           actorRequest `json:"actor"`
	RecipientName   string       `json:"recipientName"`
	RecipientEmail  string       `json:"recipientEmail"`
	Subject         string       `json:"subject"`
	Body            string       `json:"body"`
	RelatedNodeHIDs []string     `json:"relatedNodeHids"`
}

type replyMessageRequest struct {
	Actor   actorRequest `json:"actor"`
	Subject string       `json:"subject"`
	Body    string       `json:"body"`
}

func (api api) listMessagesHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	email, err := emailFromRequest(request)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	page, err := parsePagination(request)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	sortKey, err := messageSortKey(request.URL.Query().Get("sort"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	direction, err := sortDirection(request.URL.Query().Get("direction"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	result, err := api.listMessages(request.Context(), email, page, sortKey, direction)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, result)
}

func (api api) getMessageHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	message, err := api.findMessage(request.Context(), chi.URLParam(request, "messageId"))
	if err != nil {
		handleNeo4jReadError(writer, err, "message not found")
		return
	}

	writeJSON(writer, http.StatusOK, message)
}

func (api api) unreadMessageCountHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	email, err := emailFromRequest(request)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	session := api.driver.NewSession(request.Context(), neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(request.Context())
	value, err := session.ExecuteRead(request.Context(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(request.Context(), `
MATCH (:Mailbox {MailboxID: $mailboxID})-[:HAS_MESSAGE]->(message:Message)
WHERE coalesce(message.IsRead, false) = false AND coalesce(message.IsDeleted, false) = false
RETURN count(message) AS unread
`, map[string]any{"mailboxID": email})
		if err != nil {
			return int64(0), err
		}
		record, err := result.Single(request.Context())
		if err != nil {
			return int64(0), err
		}
		countValue, _ := record.Get("unread")
		count, _ := countValue.(int64)
		return count, nil
	})
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	count, _ := value.(int64)
	writeJSON(writer, http.StatusOK, unreadCountResponse{UnreadCount: count})
}

func (api api) createMessageHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	var payload directMessageRequest
	if err := decodeJSON(request, &payload); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	actor, err := actorFromRequest(request, metadata.Actor{Name: payload.Actor.Name, Email: payload.Actor.Email, Admin: payload.Actor.Admin})
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	if payload.RecipientName == "" || payload.RecipientEmail == "" {
		writeError(writer, http.StatusBadRequest, "recipient name and email are required")
		return
	}

	messageID, err := api.appendDirectMessage(request.Context(), directMessageInput{
		Sender:          actor,
		Recipient:       metadata.Actor{Name: payload.RecipientName, Email: payload.RecipientEmail},
		Subject:         payload.Subject,
		Body:            payload.Body,
		RelatedNodeHIDs: payload.RelatedNodeHIDs,
		ReplyTo:         "",
		Now:             api.now(),
	})
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	message, err := api.findMessage(request.Context(), messageID)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusCreated, message)
}

func (api api) replyMessageHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	var payload replyMessageRequest
	if err := decodeJSON(request, &payload); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	actor, err := actorFromRequest(request, metadata.Actor{Name: payload.Actor.Name, Email: payload.Actor.Email, Admin: payload.Actor.Admin})
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	original, err := api.findMessage(request.Context(), chi.URLParam(request, "messageId"))
	if err != nil {
		handleNeo4jReadError(writer, err, "message not found")
		return
	}

	subject := payload.Subject
	if subject == "" {
		subject = "Re: " + original.Subject
	}

	messageID, err := api.appendDirectMessage(request.Context(), directMessageInput{
		Sender:          actor,
		Recipient:       metadata.Actor{Name: original.Sender, Email: original.SenderEmail},
		Subject:         subject,
		Body:            payload.Body,
		RelatedNodeHIDs: original.RelatedNodeHIDs,
		ReplyTo:         original.MessageID,
		Now:             api.now(),
	})
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	message, err := api.findMessage(request.Context(), messageID)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusCreated, message)
}

func (api api) markMessageReadHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	message, err := api.updateMessageFlags(request.Context(), chi.URLParam(request, "messageId"), map[string]any{
		"IsRead": true,
		"ReadAt": api.now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		handleNeo4jReadError(writer, err, "message not found")
		return
	}

	writeJSON(writer, http.StatusOK, message)
}

func (api api) deleteMessageHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	message, err := api.updateMessageFlags(request.Context(), chi.URLParam(request, "messageId"), map[string]any{
		"IsDeleted": true,
		"DeletedAt": api.now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		handleNeo4jReadError(writer, err, "message not found")
		return
	}

	writeJSON(writer, http.StatusOK, message)
}

type directMessageInput struct {
	Sender          metadata.Actor
	Recipient       metadata.Actor
	Subject         string
	Body            string
	RelatedNodeHIDs []string
	ReplyTo         string
	Now             time.Time
}

func (api api) listMessages(ctx context.Context, email string, page pageRequest, sortKey string, direction string) (listResponse[messageResponse], error) {
	session := api.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := fmt.Sprintf(`
MATCH (:Mailbox {MailboxID: $mailboxID})-[:HAS_MESSAGE]->(message:Message)
WHERE coalesce(message.IsDeleted, false) = false
WITH message,
     coalesce(message.Subject, "") AS subjectSort,
     coalesce(message.SentAt, "") AS datetimeSort,
     coalesce(message.Sender, "") AS senderSort,
     coalesce(message.RelatedNodeHIDs[0], "") AS hidSort
RETURN properties(message) AS properties
ORDER BY %s %s
SKIP $skip
LIMIT $limit
`, sortKey, direction)
		records, err := tx.Run(ctx, query, map[string]any{"mailboxID": email, "skip": page.Offset, "limit": page.Limit})
		if err != nil {
			return nil, err
		}
		collected, err := records.Collect(ctx)
		if err != nil {
			return nil, err
		}

		items := make([]messageResponse, 0, len(collected))
		for _, record := range collected {
			items = append(items, messageFromProperties(propertiesFromRecord(record, "properties"), false))
		}

		countResult, err := tx.Run(ctx, `
MATCH (:Mailbox {MailboxID: $mailboxID})-[:HAS_MESSAGE]->(message:Message)
WHERE coalesce(message.IsDeleted, false) = false
RETURN count(message) AS total
`, map[string]any{"mailboxID": email})
		if err != nil {
			return nil, err
		}
		countRecord, err := countResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		totalValue, _ := countRecord.Get("total")
		total, _ := totalValue.(int64)
		return listResponse[messageResponse]{Items: items, Page: page.Page, Limit: page.Limit, Total: total}, nil
	})
	if err != nil {
		return listResponse[messageResponse]{}, err
	}
	response, _ := result.(listResponse[messageResponse])
	return response, nil
}

func (api api) findMessage(ctx context.Context, messageID string) (messageResponse, error) {
	session := api.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(ctx)

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (message:Message {MessageID: $messageID})
RETURN properties(message) AS properties
LIMIT 1
`, map[string]any{"messageID": messageID})
		if err != nil {
			return messageResponse{}, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return messageResponse{}, err
		}
		return messageFromProperties(propertiesFromRecord(record, "properties"), true), nil
	})
	if err != nil {
		return messageResponse{}, err
	}
	message, ok := value.(messageResponse)
	if !ok {
		return messageResponse{}, errors.New("unexpected message response")
	}
	return message, nil
}

func (api api) appendDirectMessage(ctx context.Context, input directMessageInput) (string, error) {
	session := api.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(ctx)

	value, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		messageID := identity.NewUUID()
		sentAt := input.Now.UTC().Format(time.RFC3339)
		result, err := tx.Run(ctx, `
MERGE (user:User {UserEmail: $recipientEmail})
SET user.UserName = $recipientName,
    user.UserHash = $recipientEmail
MERGE (user)-[:OWNS_MAILBOX]->(mailbox:Mailbox {MailboxID: $recipientEmail})
ON CREATE SET mailbox.Owner = $recipientName,
              mailbox.OwnerEmail = $recipientEmail,
              mailbox.UnreadCount = 0,
              mailbox.Created = $sentAt
SET mailbox.LastTouch = $sentAt
CREATE (message:Message {
  MessageID: $messageID,
  uuid: $messageID,
  Subject: $subject,
  Body: $body,
  MessageType: $messageType,
  SentAt: $sentAt,
  ReadAt: "Null",
  DeletedAt: "Null",
  Sender: $sender,
  SenderEmail: $senderEmail,
  Recipient: $recipientName,
  RecipientEmail: $recipientEmail,
  RelatedNodeHIDs: $relatedNodeHIDs,
  RelatedRelationshipTypes: [],
  CommitID: "Null",
  ReplyToMessageID: $replyTo,
  ChangeTypeSummary: "direct message",
  OldOwner: "Null",
  CurrentOwner: "Null",
  IsRead: false,
  IsDeleted: false,
  RequiresApproval: false,
  ApprovalStatus: "NOT_APPLICABLE"
})
CREATE (mailbox)-[:HAS_MESSAGE]->(message)
SET mailbox.UnreadCount = coalesce(mailbox.UnreadCount, 0) + 1
RETURN message.MessageID AS messageID
`, map[string]any{
			"recipientName":   input.Recipient.Name,
			"recipientEmail":  input.Recipient.Email,
			"messageID":       messageID,
			"subject":         input.Subject,
			"body":            input.Body,
			"messageType":     string(messaging.MessageTypeDirect),
			"sentAt":          sentAt,
			"sender":          input.Sender.Name,
			"senderEmail":     input.Sender.Email,
			"relatedNodeHIDs": input.RelatedNodeHIDs,
			"replyTo":         nullWhenEmpty(input.ReplyTo),
		})
		if err != nil {
			return "", err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return "", err
		}
		value, _ := record.Get("messageID")
		return stringValue(value), nil
	})
	if err != nil {
		return "", err
	}
	messageID, _ := value.(string)
	if messageID == "" {
		return "", errors.New("message creation did not return MessageID")
	}
	return messageID, nil
}

func (api api) updateMessageFlags(ctx context.Context, messageID string, props map[string]any) (messageResponse, error) {
	session := api.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(ctx)

	value, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (message:Message {MessageID: $messageID})
SET message += $props
RETURN properties(message) AS properties
`, map[string]any{"messageID": messageID, "props": props})
		if err != nil {
			return messageResponse{}, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return messageResponse{}, err
		}
		return messageFromProperties(propertiesFromRecord(record, "properties"), true), nil
	})
	if err != nil {
		return messageResponse{}, err
	}
	message, _ := value.(messageResponse)
	return message, nil
}

func messageSortKey(value string) (string, error) {
	switch value {
	case "", "datetime":
		return "datetimeSort", nil
	case "subject":
		return "subjectSort", nil
	case "hid":
		return "hidSort", nil
	case "sender":
		return "senderSort", nil
	default:
		return "", fmt.Errorf("unsupported message sort %q", value)
	}
}

func messageFromProperties(properties map[string]any, includeBody bool) messageResponse {
	body := ""
	if includeBody {
		body = stringProperty(properties, "Body")
	}
	isRead, _ := properties["IsRead"].(bool)
	return messageResponse{
		MessageID:                stringProperty(properties, "MessageID"),
		Subject:                  stringProperty(properties, "Subject"),
		Body:                     body,
		SentAt:                   stringProperty(properties, "SentAt"),
		Sender:                   stringProperty(properties, "Sender"),
		SenderEmail:              stringProperty(properties, "SenderEmail"),
		MessageType:              stringProperty(properties, "MessageType"),
		IsRead:                   isRead,
		RelatedNodeHIDs:          anySliceToStrings(properties["RelatedNodeHIDs"]),
		RelatedRelationshipTypes: anySliceToStrings(properties["RelatedRelationshipTypes"]),
		Properties:               properties,
	}
}

func nullWhenEmpty(value string) string {
	if value == "" {
		return metadata.NullValue
	}

	return value
}
