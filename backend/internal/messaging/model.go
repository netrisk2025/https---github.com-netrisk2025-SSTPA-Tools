// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/metadata"
)

type MessageType string

const (
	MessageTypeDirect             MessageType = "DIRECT"
	MessageTypeChangeNotification MessageType = "CHANGE_NOTIFICATION"
	MessageTypeSystem             MessageType = "SYSTEM"
)

type ApprovalStatus string

const (
	ApprovalNotApplicable ApprovalStatus = "NOT_APPLICABLE"
	ApprovalPending       ApprovalStatus = "PENDING"
	ApprovalApproved      ApprovalStatus = "APPROVED"
	ApprovalRejected      ApprovalStatus = "REJECTED"
)

type ChangeNotification struct {
	CommitID                 string
	Subject                  string
	Body                     string
	SentAt                   time.Time
	Sender                   metadata.Actor
	Recipient                metadata.Actor
	RelatedNodeHIDs          []string
	RelatedRelationshipTypes []string
	ChangeTypeSummary        string
	OldOwner                 string
	CurrentOwner             string
}

func AppendChangeNotification(ctx context.Context, tx neo4j.ManagedTransaction, notification ChangeNotification) (string, error) {
	if notification.Recipient.Name == "" || notification.Recipient.Email == "" {
		return "", fmt.Errorf("recipient name and email are required")
	}

	if notification.Sender.Name == "" || notification.Sender.Email == "" {
		return "", fmt.Errorf("sender name and email are required")
	}

	messageID := identity.NewUUID()
	mailboxID := notification.Recipient.Email
	subject := notification.Subject
	if subject == "" {
		subject = "SSTPA data changed"
	}

	params := map[string]any{
		"userName":                 notification.Recipient.Name,
		"userEmail":                notification.Recipient.Email,
		"userHash":                 notification.Recipient.Email,
		"mailboxID":                mailboxID,
		"messageID":                messageID,
		"subject":                  subject,
		"body":                     notification.Body,
		"messageType":              string(MessageTypeChangeNotification),
		"sentAt":                   notification.SentAt.UTC().Format(time.RFC3339),
		"sender":                   notification.Sender.Name,
		"senderEmail":              notification.Sender.Email,
		"recipient":                notification.Recipient.Name,
		"recipientEmail":           notification.Recipient.Email,
		"relatedNodeHIDs":          notification.RelatedNodeHIDs,
		"relatedRelationshipTypes": notification.RelatedRelationshipTypes,
		"commitID":                 notification.CommitID,
		"changeTypeSummary":        notification.ChangeTypeSummary,
		"oldOwner":                 notification.OldOwner,
		"currentOwner":             notification.CurrentOwner,
	}

	result, err := tx.Run(ctx, `
MERGE (u:User {UserEmail: $userEmail})
SET u.UserName = $userName,
    u.UserHash = $userHash
MERGE (u)-[:OWNS_MAILBOX]->(mailbox:Mailbox {MailboxID: $mailboxID})
ON CREATE SET mailbox.Owner = $userName,
              mailbox.OwnerEmail = $userEmail,
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
  Recipient: $recipient,
  RecipientEmail: $recipientEmail,
  RelatedNodeHIDs: $relatedNodeHIDs,
  RelatedRelationshipTypes: $relatedRelationshipTypes,
  CommitID: $commitID,
  ChangeTypeSummary: $changeTypeSummary,
  OldOwner: $oldOwner,
  CurrentOwner: $currentOwner,
  IsRead: false,
  IsDeleted: false,
  RequiresApproval: false,
  ApprovalStatus: "NOT_APPLICABLE"
})
CREATE (mailbox)-[:HAS_MESSAGE]->(message)
SET mailbox.UnreadCount = coalesce(mailbox.UnreadCount, 0) + 1
RETURN message.MessageID AS messageID
`, params)
	if err != nil {
		return "", err
	}

	record, err := result.Single(ctx)
	if err != nil {
		return "", err
	}

	value, _ := record.Get("messageID")
	got, ok := value.(string)
	if !ok || got == "" {
		return "", fmt.Errorf("message creation did not return MessageID")
	}

	return got, nil
}
