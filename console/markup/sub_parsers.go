package markup

import (
	"fmt"
	"strconv"
	"strings"
)

// Control Sequence Introducer.
// Provides in-band signalling to a terminal that styling, font, colour or location data is following.
const csi = "\033["

var formatters = map[string]string{
	"/":      "0m",
	"bold":   "1m",
	"strong": "1m",
	"italic": "3m",
	// TODO: Add support for underline styles and colours.
	"underline": "4m",
	"blink":     "5m",
	"invert":    "7m",
	"strike":    "9m",
	"strikeout": "9m",

	// foreground colours
	"black":          "30m",
	"red":            "31m",
	"green":          "32m",
	"yellow":         "33m",
	"blue":           "34m",
	"magenta":        "35m",
	"cyan":           "36m",
	"white":          "37m",
	"gray":           "90m",
	"bright-black":   "90m",
	"bright-red":     "91m",
	"bright-green":   "92m",
	"bright-yellow":  "93m",
	"bright-blue":    "94m",
	"bright-magenta": "95m",
	"bright-cyan":    "96m",
	"bright-white":   "97m",

	// background colours
	"bg-black":          "40m",
	"bg-red":            "41m",
	"bg-green":          "42m",
	"bg-yellow":         "43m",
	"bg-blue":           "44m",
	"bg-magenta":        "45m",
	"bg-cyan":           "46m",
	"bg-white":          "47m",
	"bg-gray":           "40m",
	"bg-bright-black":   "100m",
	"bg-bright-red":     "101m",
	"bg-bright-green":   "102m",
	"bg-bright-yellow":  "103m",
	"bg-bright-blue":    "104m",
	"bg-bright-magenta": "105m",
	"bg-bright-cyan":    "106m",
	"bg-bright-white":   "107m",
}

// TODO: Implement.
// We support https://no-color.org/
// var noColour bool

// func init() {
// 	noColour = slices.ContainsFunc(os.Environ(), func(envVar string) bool {
// 		return envVar == "NO_COLOUR" || envVar == "NO_COLOR"
// 	})
// }

func parseTag(token token) (value string, ok bool) {
	elements := parseTagElements(strings.ToLower(token.value))
	tag := elements[0]

	if len(elements) == 1 {
		if v, ok := formatters[tag]; ok {
			return fmt.Sprintf("%v%v", csi, v), ok
		}
	}

	switch {
	case tag == "rgb":
		if s, ok := parseRgb(elements[1:]); ok {
			return s, true
		}

	case tag == "hex":
		if len(elements) == 2 {
			if s, ok := parseHex(elements[1]); ok {
				return s, true
			}
		}

	case tag == "link":
		if len(elements) == 3 {
			// https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda
			// \e]8;;http://example.com\e\\This is a link\e]8;;\e\\\n
			// $"{ESC}]8;id={linkId};{link}{ESC}\\{ansi}{ESC}]8;;{ESC}\\";
			return fmt.Sprintf("\x1b]8;id=%v;%v\x1b\\%v\x1b]8;;\x1b\\", elements[2], elements[2], elements[1]), true
		}
	}

	// We tried.  But this isn't tag.
	return token.value, false
}

func parseTagElements(parameters string) []string {
	if len(parameters) > 2 {
		// Parameters format: "<name::param1::param2::...|paramN>"
		// We trim the opening and closing tag.
		return strings.Split(parameters[1:len(parameters)-1], "::")
	}

	return []string{}
}

func parseHex(hexTriplet string) (value string, ok bool) {
	// We support the optional leading #.
	if hexTriplet[0] == '#' {
		hexTriplet = hexTriplet[1:]
	}

	if len(hexTriplet) != 6 {
		return "", false
	}

	colours := [3]int64{0, 0, 0}
	for i, v := range []string{hexTriplet[0:2], hexTriplet[2:4], hexTriplet[4:6]} {
		i64, err := strconv.ParseInt(v, 16, 64)
		if err != nil {
			return "", false
		}

		colours[i] = i64
	}

	return fmt.Sprintf("%v38;2;%v;%v;%vm", csi, colours[0], colours[1], colours[2]), true
}

func parseRgb(rgb []string) (value string, ok bool) {
	if len(rgb) != 3 {
		return "", false
	}

	colours := [3]int64{}
	for i, parameter := range rgb {
		value, err := strconv.ParseInt(parameter, 10, 64)
		if err != nil {
			return "", false
		}
		colours[i] = value
	}

	return fmt.Sprintf("%v38;2;%v;%v;%vm", csi, colours[0], colours[1], colours[2]), true
}
