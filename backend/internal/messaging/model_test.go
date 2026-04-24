// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package messaging

import "testing"

func TestMessageTypeValuesMatchSRS(t *testing.T) {
	if MessageTypeDirect != "DIRECT" {
		t.Fatalf("DIRECT enum changed: %q", MessageTypeDirect)
	}
	if MessageTypeChangeNotification != "CHANGE_NOTIFICATION" {
		t.Fatalf("CHANGE_NOTIFICATION enum changed: %q", MessageTypeChangeNotification)
	}
	if MessageTypeSystem != "SYSTEM" {
		t.Fatalf("SYSTEM enum changed: %q", MessageTypeSystem)
	}
}
