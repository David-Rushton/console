// Control your console.
package console

import "fmt"

// Control Sequence Introducer
const csi = "\033["

// Clears the console.
func Clear() {
	Position(1, 1)
	fmt.Printf("%v2J", csi)
}

// Positions the console cursor.
// The top left corner is 1, 1.
func Position(top, left int) {
	fmt.Printf("%v%v;%vH", csi, top, left)
}
