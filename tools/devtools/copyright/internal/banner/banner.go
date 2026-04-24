package banner

import "strings"

const Marker = "SSTPA Tools software and all associated modules"

const Text = `// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
`

func HasBanner(content string) bool {
	return strings.Contains(content, Marker)
}

func Prepend(content string) string {
	if HasBanner(content) {
		return content
	}

	if strings.HasPrefix(content, "#!") {
		idx := strings.Index(content, "\n")
		if idx < 0 {
			return content + "\n" + Text
		}

		return content[:idx+1] + Text + content[idx+1:]
	}

	return Text + content
}
