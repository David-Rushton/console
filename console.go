/*
Interactive terminal commands.

This package wraps the [ANSI escape code] Control Sequence Introducer commands.  Also known as CSI
commands.

Requires [ANSI escape code] support.

[ANSI escape code]: https://en.wikipedia.org/wiki/ANSI_escape_code#Control_Sequence_Introducer_commands
*/
package console

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

// Rings the terminal bell.
//
// Bell is not supported by all terminals.  Some terminals allow users to disable the bell.  When
// this is the case calls to bell are effectively a no-op.
func Bell() {
	fmt.Print("\x07")
}

// Clears the entire screen.
//
// The cursor is moved to the top left corner.
func Clear() {
	fmt.Print("\x1b[2J")

	// Above command moves the cursor in some terminals but not others.
	// Resetting the position ensures consistent behaviour.
	SetCursorPosition(1, 1)
}

// Shows the terminal cursor.
func ShowCursor() {
	fmt.Print("\x1b[?25h")
}

// Hides the terminal cursor.
func HideCursor() {
	fmt.Print("\x1b[?25l")
}

// Moves the cursor up n rows.
//
// If the cursor reaches the top of the screen it wil stop moving.  The terminal will not scroll.
func CursorUp(n uint) {
	fmt.Printf("\x1b[%dA", n)
}

// Moves the cursor down n rows.
//
// If the cursor reaches the bottom of the screen it wil stop moving.  The terminal will not scroll.
func MoveCursorDown(n uint) {
	fmt.Printf("\x1b[%dB", n)
}

// Moves the cursor forward n columns.
//
// If the cursor reaches the edge of the screen it wil stop moving.  It will not wrap to the next
// line.
func MoveCursorForward(n uint) {
	fmt.Printf("\x1b[%dC", n)
}

// Moves the cursor forward n columns.
//
// If the cursor reaches the edge of the screen it wil stop moving.  It will not wrap to the
// previous line.
func MoveCursorBack(n uint) {
	fmt.Printf("\x1b[%dD", n)
}

// Moves the cursor to the start of the line n rows below current.
//
// If the cursor reaches the bottom of the screen it will stop.  The terminal will not scroll.
func MoveCursorNextLine(n uint) {
	fmt.Printf("\x1b[%dE", n)
}

// Moves the cursor to the start of the line n rows above current.
//
// If the cursor reaches the top of the screen it will stop.  The terminal will not scroll.
func MoveCursorPreviousLine(n uint) {
	fmt.Printf("\x1b[%dF", n)
}

// Returns the cursor position.
//
// Position is relative to top left.  Which is 1, 1.
func GetCursorPosition() (column, row uint, err error) {
	cursorPosition, err := runInRawMode("\x1b[6n")
	if err != nil {
		return 0, 0, fmt.Errorf("cannot query for cursor position because %v", err)
	}

	var r, c uint
	_, err = fmt.Sscanf(cursorPosition, "\x1b[%d;%dR", &r, &c)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot find cursor position within response: %v", cursorPosition)
	}

	return r, c, nil
}

// Moves the cursor to the requested row and column.
//
// The top left of the screen is always 1, 1.  You cannot move the cursor to a rows that have
// scrolled off the screen.  If row or column are 0 then 1 is assumed.
func SetCursorPosition(column, row uint) {
	if column == 0 {
		column = 1
	}

	if row == 0 {
		row = 1
	}

	fmt.Printf("\x1b[%d;%dH", row, column)
}

// Returns the visible size of the terminal.
func GetSize() (columns, rows int, err error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

func runInRawMode(command string) (result string, err error) {
	// Raw mode executes commands in the terminal:
	//  - Disables echoing
	//  - Disables line buffering
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Write to the terminal.
	if _, err = os.Stdout.WriteString(command); err != nil {
		return "", err
	}

	// Read the response.
	// We need to handle the case where the response never arrives.
	select {
	case <-time.After(time.Millisecond * 200):
		return "", fmt.Errorf("standard in timed out")
	default:
		var buf [256]byte
		n, err := os.Stdin.Read(buf[:])
		if err != nil {
			return "", err
		}

		return string(buf[:n]), nil
	}
}
