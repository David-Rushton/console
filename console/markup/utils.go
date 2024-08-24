package markup

// import (
// 	"strconv"
// 	"strings"
// )

// // Converts a hex triplet to a array of int.
// // Where red = colours[0], green = colours[1] & blue = colours[2].
// func parseRgb(hexTriplet string) (colours [3]int64, ok bool) {
// 	// remove optional # prefix
// 	if strings.HasPrefix(hexTriplet, "#") {
// 		hexTriplet = hexTriplet[1:]
// 	}

// 	if len(hexTriplet) != 6 {
// 		return colours, false
// 	}

// 	for i, v := range [3]int64{0, 2, 4} {
// 		value, err := strconv.ParseInt(hexTriplet[v:v+2], 16, 64)
// 		if err != nil {
// 			return colours, false
// 		}

// 		colours[i] = value
// 	}

// 	return colours, true
// }
