// Markup is converted to Ansi Escape Codes, providing rich control
// over your console output.
package markup

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

var debugMode bool

func init() {
	debugMode = slices.Contains(os.Args, "--markup-debug")
}

// Converts markup stings into a rich console output.
func Parse(markup string) string {
	buffer := []string{}

	for token := range tokenise(markup) {
		// TODO: Does standard-out make sense for a CLI app?
		if debugMode {
			fmt.Printf("%#v\n", token)
		}

		if token.tokenType == tag {
			if value, ok := parseTag(token); ok {
				buffer = append(buffer, value)
				continue
			}
		}

		buffer = append(buffer, token.value)
	}

	return strings.Join(buffer, "")
}
