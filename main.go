package main

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/gbataille/zap_log_prettier/consolefmt"
	"github.com/pterm/pterm"
)

func main() {
	handleSTDIN()
}

func handleSTDIN() {
	// Enable debug messages.
	pterm.EnableDebugMessages()

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		lineB := input.Bytes()
		logLine, err := fromJsonGeneric(lineB)
		if err != nil {
			pterm.Debug.Println(string(lineB))
		} else {
			toHumanLog(logLine)
		}
	}
}

func fromJsonGeneric(in []byte) (map[string]any, error) {
	jsonMap := make(map[string](any))
	err := json.Unmarshal(in, &jsonMap)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}

func toHumanLog(logLine map[string]any) {
	methodFromLine(logLine)(logLine)
}

func methodFromLine(logLine map[string]any) func(a ...any) {
	levelRaw, ok := logLine["level"]

	if !ok {
		return pterm.Println
	}

	level, ok := levelRaw.(string)
	if !ok {
		return pterm.Println
	}

	return methodFromLevel(level)
}

func methodFromLevel(level string) func(a ...any) {
	switch strings.ToUpper(level) {
	case "ERROR":
		return withNoReturn(consolefmt.Error.Println)
	case "INFO":
		return withNoReturn(consolefmt.Info.Println)
	case "FATAL":
		return withNoReturn(consolefmt.Fatal.Println)
	case "DEBUG":
		return withNoReturn(consolefmt.Debug.Println)
	case "WARNING":
		return withNoReturn(consolefmt.Warning.Println)
	case "WARN":
		return withNoReturn(consolefmt.Warning.Println)
	default:
		return pterm.Println
	}
}

func withNoReturn(f func(a ...any) *pterm.TextPrinter) func(a ...any) {
	return func(a ...any) {
		f(a)
	}
}
