package raw

import (
	"fmt"

	"github.com/pterm/pterm"
)

func HandleRAWLine(line string) {
	pterm.ThemeDefault.PrimaryStyle.Println(line)
	fmt.Println()
}
