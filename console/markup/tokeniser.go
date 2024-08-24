package markup

import (
	"iter"
	"strings"
)

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
