package raw

import (
	"github.com/pterm/pterm"
)

func HandleRAWLine(line string) {
	pterm.ThemeDefault.PrimaryStyle.Println(line)
}
