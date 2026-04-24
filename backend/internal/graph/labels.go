// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package graph

import "sstpa-tool/backend/internal/identity"

func LabelFor(nodeType identity.NodeType) (string, bool) {
	if _, ok := identity.TypeID(nodeType); !ok {
		return "", false
	}

	return string(nodeType), true
}
