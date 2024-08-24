// Markup is converted to Ansi Escape Codes, providing rich control
// over your console output.
package markup

import (
	"fmt"
	"iter"
	"log"
	"strconv"
	"strings"
)

// Converts markup stings into a rich console output.
func Parse(markup string) string {
	buffer := []string{}

	for token := range tokenise(markup) {
		// DEBUG: Remove later.
		// fmt.Printf("%#v\n", token)

		if token.tokenType == tag {
			buffer = append(buffer, convertTag(token))
			continue
		}

		buffer = append(buffer, token.value)
	}

	return strings.Join(buffer, "")
}

// Control Sequence Introducer
const csi = "\033["

// Converts a markup tag into a ANSI escaped string.
func convertTag(token token) string {
	tag := strings.ToLower(token.value)

	switch {
	case strings.HasPrefix(tag, "<rgb"):
		if s, ok := ReadRgbTag(token); ok {
			return s
		}

	case tag == "</>":
		return fmt.Sprintf("%v0m", csi)
	}

	// We tried.  But this isn't tag.
	return token.value
}

func ReadRgbTag(token token) (value string, ok bool) {
	elements := strings.Split(token.value, ":")

	if len(elements) != 5 {
		return "", false
	}

	colours := [3]int64{}
	for i := 1; i < 4; i++ {
		value, err := strconv.ParseInt(elements[i], 10, 64)
		if err != nil {
			log.Fatal(err)
			return "", false
		}
		colours[i-1] = value
	}

	return fmt.Sprintf("%v38;2;%v;%v;%vm", csi, colours[0], colours[1], colours[2]), true
}

// Tokens are assigned a type, to simplify parsing.
type tokenType int

const (
	whitespace tokenType = iota
	tag
	literal
)

// Describes a word.
type token struct {
	tokenType tokenType
	value     string
}

// Converts a markup string into a series of tokens.
func tokenise(markup string) iter.Seq[token] {
	return func(yield func(token) bool) {
		const openTag = '<'
		const closeTag = '>'
		const space = ' '

		buffer := []rune{}
		flush := func() {
			if len(buffer) == 0 {
				return
			}

			value := string(buffer)
			tokenType := literal

			if strings.HasPrefix(value, string(openTag)) {
				if strings.HasSuffix(value, string(closeTag)) {
					tokenType = tag
				}
			}

			buffer = []rune{}

			yield(token{tokenType, value})
		}

		for _, r := range markup {
			switch r {
			// Opening tags always start a new buffer, to simplify parsing.
			case openTag:
				flush()
				buffer = append(buffer, r)

			// Closing tags always end a buffer, as this simplifies parsing.
			case closeTag:
				buffer = append(buffer, r)
				flush()

			// Always flush on these whitespace.
			// Tokens that exclusively contain whitespace simplifies parsing.
			case space, '\t', '\r', '\n':
				flush()
				yield(token{whitespace, string(r)})

			default:
				buffer = append(buffer, r)
			}
		}

		// Make sure the buffer is empty before we exit.
		flush()
	}
}
