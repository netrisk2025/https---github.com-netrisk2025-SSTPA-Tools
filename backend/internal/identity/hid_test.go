// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package identity

import "testing"

func TestFormatHID(t *testing.T) {
	tests := []struct {
		name     string
		typeID   string
		index    string
		sequence int
		want     string
		wantErr  bool
	}{
		{"capability", "CAP", "", 0, "CAP__0", false},
		{"root system", "SYS", "1", 0, "SYS_1_0", false},
		{"nested element", "EL", "1.2.3", 4, "EL_1.2.3_4", false},
		{"negative sequence", "SYS", "1", -1, "", true},
		{"unknown type", "XYZ", "1", 0, "", true},
		{"malformed index", "SYS", "1..2", 0, "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := FormatHID(test.typeID, test.index, test.sequence)
			if (err != nil) != test.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, test.wantErr)
			}

			if got != test.want {
				t.Fatalf("got %q want %q", got, test.want)
			}
		})
	}
}

func TestParseHIDRoundTrip(t *testing.T) {
	in := "EL_1.2.3_4"
	typeID, index, sequence, err := ParseHID(in)
	if err != nil {
		t.Fatal(err)
	}

	got, err := FormatHID(typeID, index, sequence)
	if err != nil || got != in {
		t.Fatalf("round-trip failed: got %q err %v", got, err)
	}
}

func TestParseHIDRejectsMalformed(t *testing.T) {
	for _, bad := range []string{"", "NOPE", "SYS_1_", "SYS__", "SYS_1_abc", "SYS_1_0_extra"} {
		if _, _, _, err := ParseHID(bad); err == nil {
			t.Errorf("expected error for %q", bad)
		}
	}
}
