// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package manifest

import "testing"

func TestNormalizedReferenceItemValidate(t *testing.T) {
	item := NormalizedReferenceItem{
		FrameworkName:    "MITRE ATT&CK",
		FrameworkVersion: "16.1",
		ExternalID:       "T1001",
		ExternalType:     "Technique",
		Name:             "Data Obfuscation",
		ShortDescription: "Example entry",
		LongDescription:  "Null",
		SourceURI:        "offline://attack/16.1/T1001",
	}

	if err := item.Validate(); err != nil {
		t.Fatalf("expected item to validate: %v", err)
	}
}

func TestNormalizedReferenceItemValidateRejectsMissingFields(t *testing.T) {
	item := NormalizedReferenceItem{}

	if err := item.Validate(); err == nil {
		t.Fatal("expected validation error for empty item")
	}
}
