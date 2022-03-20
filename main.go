package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
			fmt.Println(string(lineB))
		}
		toHumanLog(logLine)
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
		return withNoReturn(pterm.Error.Println)
	case "INFO":
		return withNoReturn(pterm.Info.Println)
	case "FATAL":
		return withNoReturn(pterm.Info.Println)
	case "DEBUG":
		return withNoReturn(pterm.Debug.Println)
	case "WARNING":
		return withNoReturn(pterm.Debug.Println)
	default:
		return pterm.Println
	}
}

func withNoReturn(f func(a ...any) *pterm.TextPrinter) func(a ...any) {
	return func(a ...any) {
		f(a)
	}
}
