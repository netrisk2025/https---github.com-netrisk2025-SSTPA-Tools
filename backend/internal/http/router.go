// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package apihttp

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/onboarding"
)

type RouterOptions struct {
	Version      string
	Driver       neo4j.DriverWithContext
	DatabaseName string
}

func NewRouter(version string) http.Handler {
	return NewRouterWithOptions(RouterOptions{Version: version})
}

func NewRouterWithOptions(options RouterOptions) http.Handler {
	api := newAPI(options)
	router := chi.NewRouter()
	router.Get("/healthz", healthHandler(api.version))
	router.Route("/api/v1", func(group chi.Router) {
		group.Get("/health", healthHandler(api.version))
		group.Get("/openapi.yaml", api.openapiHandler)

		group.Get("/nodes", api.listNodesHandler)
		group.Get("/nodes/uuid/{uuid}", api.getNodeByUUIDHandler)
		group.Get("/nodes/{hid}", api.getNodeByHIDHandler)
		group.Get("/nodes/{hid}/context", api.nodeContextHandler)
		group.Get("/hierarchy", api.hierarchyHandler)
		group.Get("/search", api.searchHandler)
		group.Post("/validate/relationship", api.validateRelationshipHandler)
		group.Post("/mutations", api.mutationsHandler)

		group.Get("/messages/unread-count", api.unreadMessageCountHandler)
		group.Get("/messages", api.listMessagesHandler)
		group.Post("/messages", api.createMessageHandler)
		group.Get("/messages/{messageId}", api.getMessageHandler)
		group.Post("/messages/{messageId}/reply", api.replyMessageHandler)
		group.Post("/messages/{messageId}/read", api.markMessageReadHandler)
		group.Delete("/messages/{messageId}", api.deleteMessageHandler)

		group.Get("/reference/frameworks", api.listReferenceFrameworksHandler)
		group.Get("/reference/items", api.listReferenceItemsHandler)
		group.Get("/reference/items/uuid/{uuid}", api.getReferenceItemByUUIDHandler)
		group.Get("/reference/items/{externalID}", api.getReferenceItemByExternalIDHandler)
		group.Get("/reference/items/{uuid}/related", api.relatedReferenceItemsHandler)
		group.Get("/reference/search", api.searchReferenceItemsHandler)
		group.Post("/reference/validate-assignment", api.validateReferenceAssignmentHandler)

		group.Get("/references/assignments/{sourceHID}", api.listReferenceAssignmentsHandler)
		group.Post("/references/assignments", api.createReferenceAssignmentHandler)
		group.Delete("/references/assignments", api.deleteReferenceAssignmentHandler)
		group.Get("/users", api.listOnboardingHandler(onboarding.UserKind))
		group.Post("/users", api.createOnboardingHandler(onboarding.UserKind))
		group.Get("/users/{uuid}", api.getOnboardingHandler(onboarding.UserKind))
		group.Get("/admins", api.listOnboardingHandler(onboarding.AdminKind))
		group.Post("/admins", api.createOnboardingHandler(onboarding.AdminKind))
		group.Get("/admins/{uuid}", api.getOnboardingHandler(onboarding.AdminKind))
	})

	return router
}
