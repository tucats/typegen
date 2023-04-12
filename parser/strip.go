package parser

import (
	"fmt"
	"strings"
)

// Strip out any comments in the JSON code. Comments are not standard, but are
// placed in JSON by some apps (for example, VS Code). However, comments are not
// supported in the Go decoder package. This function locates the comments and
// removes them from the JSON string before starting the parse.
func strip(data []byte, debug bool) []byte {
	// Convert back to a string
	str := string(data)
	modified := false

	// Convert to lines of text based on break
	lines := strings.Split(str, "\n")

	// Scan over each line and search for comments
	for i := 0; i < len(lines); i++ {
		text := lines[i]

		quote := false
		escape := false
		comment := false

		for index, ch := range text {
			if escape {
				continue
			}

			if ch == '\\' {
				escape = true

				continue
			}

			if ch == '"' {
				quote = !quote

				continue
			}

			if quote {
				continue
			}

			if ch == '/' {
				if comment {
					modified = true

					if index == 0 {
						lines = append(lines[:i], lines[i+1:]...)
						i--
					} else {
						if strings.TrimSpace(text[:index-1]) == "" {
							lines = append(lines[:i], lines[i+1:]...)
							i--
						} else {
							lines[i] = text[:index-1]
						}
					}

					break
				}

				comment = true
			}
		}
	}

	str = strings.Join(lines, "\n")

	if debug && modified {
		fmt.Println("Stripped JSON:")
		fmt.Println(str)
	}

	return []byte(str)
}
