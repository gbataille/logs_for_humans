package pterm

import (
	"strings"

	"github.com/pterm/pterm"
)

// Redefines the printer for the prefix to all have the same length, for alignment purpose
var Info, Warning, Success, Error, Fatal, Debug pterm.PrefixPrinter

func init() {
	Info = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.InfoMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.InfoPrefixStyle,
			Text:  " INFO  ",
		},
	}

	// Warning returns a PrefixPrinter, which can be used to print text with a "warning" Prefix.
	Warning = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.WarningMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.WarningPrefixStyle,
			Text:  " WARN  ",
		},
	}

	// Error returns a PrefixPrinter, which can be used to print text with an "error" Prefix.
	Error = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.ErrorMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.ErrorPrefixStyle,
			Text:  " ERROR ",
		},
	}

	// Fatal returns a PrefixPrinter, which can be used to print text with an "fatal" Prefix.
	// NOTICE: Fatal terminates the application immediately! I remove the fatal: true flag from the common definition
	Fatal = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.FatalMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.FatalPrefixStyle,
			Text:  " FATAL ",
		},
	}

	// Debug Prints debug messages. By default it will only print if PrintDebugMessages is true.
	// You can change PrintDebugMessages with EnableDebugMessages and DisableDebugMessages, or by setting the variable itself.
	Debug = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.DebugMessageStyle,
		Prefix: pterm.Prefix{
			Text:  " DEBUG ",
			Style: &pterm.ThemeDefault.DebugPrefixStyle,
		},
		Debugger: true,
	}
}

func PrintWithLevel(level string, a ...interface{}) {
	switch strings.ToUpper(level) {
	case "ERROR":
		Error.Print(a)
	case "INFO":
		Info.Print(a)
	case "FATAL":
		Fatal.Print(a)
	case "DEBUG":
		Debug.Print(a)
	case "WARNING":
		Warning.Print(a)
	case "WARN":
		Warning.Print(a)
	default:
		pterm.ThemeDefault.PrimaryStyle.Print(a)
	}
}

func PrintlnWithLevel(level string, a ...interface{}) {
	PrintWithLevel(level, a)
	PrintWithLevel("\n")
}

func MethodFromLevel(level string) func(a ...interface{}) {
	switch strings.ToUpper(level) {
	case "ERROR":
		return withNoReturn(Error.Println)
	case "INFO":
		return withNoReturn(Info.Println)
	case "FATAL":
		return withNoReturn(Fatal.Println)
	case "DEBUG":
		return withNoReturn(Debug.Println)
	case "WARNING":
		return withNoReturn(Warning.Println)
	case "WARN":
		return withNoReturn(Warning.Println)
	default:
		return pterm.ThemeDefault.PrimaryStyle.Println
	}
}

func withNoReturn(f func(a ...interface{}) *pterm.TextPrinter) func(a ...interface{}) {
	return func(a ...interface{}) {
		f(a...)
	}
}
