package consolefmt

import "github.com/pterm/pterm"

// Redefines the printer for the prefix to have all the same length, for alignment purpose
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
	// NOTICE: Fatal terminates the application immediately!
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
