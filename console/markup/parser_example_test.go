package markup_test

import (
	"fmt"

	"example.com/console/console/markup"
)

func ExampleParser() {
	formattedConsoleString := markup.Parse("This text contains a <red>red section</> and a <bold>bold section</>.")

	fmt.Println(formattedConsoleString)
	// Output:
	// This text contains a \x1b[31mred section\x1b[0m and a \x1b[1mbold section\x1b[0m.
}
