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

		// Some tags can contain spaces, and other characters, that we would normally yield on.
		var buildingTag = false
		var skipNext = false

		for i, r := range markup {

			// Current rune is part of escape.
			// It has already been processed.
			if skipNext {
				skipNext = false
				continue
			}

			// Handle escapes.
			if r == openTag || r == closeTag {
				if i+1 < len(markup) {
					next := rune(markup[i+1])
					if r == next && next == openTag || next == closeTag {
						buffer = append(buffer, r)
						skipNext = true
						continue
					}
				}
			}

			switch r {
			case openTag:
				flush()
				buildingTag = true
				buffer = append(buffer, r)

			case closeTag:
				buffer = append(buffer, r)
				buildingTag = false
				flush()

			case space, '\t', '\r', '\n':
				if !buildingTag {
					flush()
					yield(token{whitespace, string(r)})
				} else {
					buffer = append(buffer, r)
				}

			default:
				buffer = append(buffer, r)
			}

		}

		// Make sure the buffer is empty before we exit.
		flush()
	}
}
