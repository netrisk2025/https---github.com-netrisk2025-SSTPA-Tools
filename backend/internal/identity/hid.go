// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package identity

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var indexPattern = regexp.MustCompile(`^(\d+(\.\d+)*)?$`)

func FormatHID(typeID string, index string, sequence int) (string, error) {
	if !IsValidTypeID(typeID) {
		return "", fmt.Errorf("unknown node type id %q", typeID)
	}

	if !indexPattern.MatchString(index) {
		return "", fmt.Errorf("invalid HID index %q", index)
	}

	if sequence < 0 {
		return "", fmt.Errorf("invalid negative HID sequence %d", sequence)
	}

	return fmt.Sprintf("%s_%s_%d", typeID, index, sequence), nil
}

func ParseHID(hid string) (string, string, int, error) {
	parts := strings.Split(hid, "_")
	if len(parts) != 3 {
		return "", "", 0, fmt.Errorf("invalid HID shape %q", hid)
	}

	typeID := parts[0]
	index := parts[1]
	sequenceValue := parts[2]

	if !IsValidTypeID(typeID) {
		return "", "", 0, fmt.Errorf("unknown node type id %q", typeID)
	}

	if !indexPattern.MatchString(index) {
		return "", "", 0, fmt.Errorf("invalid HID index %q", index)
	}

	sequence, err := strconv.Atoi(sequenceValue)
	if err != nil || sequence < 0 {
		return "", "", 0, fmt.Errorf("invalid HID sequence %q", sequenceValue)
	}

	return typeID, index, sequence, nil
}
