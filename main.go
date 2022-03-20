package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
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
	logFunc := methodFromLine(logLine)
	logMsg, ok := extract(logLine, "msg")

	if !ok {
		logFunc(fmt.Sprintf("%v", logLine))
		return
	}

	logTs, _ := extract(logLine, "ts")
	logFunc(fmt.Sprintf("[%s]", logTs), logMsg)

	caller, _ := extract(logLine, "caller")
	leveledList := pterm.LeveledList{
		pterm.LeveledListItem{Level: 0, Text: caller},
	}
	root := pterm.NewTreeFromLeveledList(leveledList)
	tree, _ := pterm.DefaultTree.WithRoot(root).Srender()
	pterm.Print(tree)

	keys := make([]string, 0, len(logLine))
	for k := range logLine {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	kv := make([][]string, 0, len(logLine))
	for _, k := range keys {
		kv = append(kv, []string{fmt.Sprintf("\t\t%s", k), asString(logLine[k])})
	}

	pterm.DefaultTable.WithData(pterm.TableData(kv)).Render()

	pterm.Println()
}

func asString(raw any) string {
	var value string
	switch raw.(type) {
	case string:
		value = raw.(string)
	case fmt.Stringer:
		value = raw.(fmt.Stringer).String()
	default:
		value = fmt.Sprintf("%v", raw)

	}
	return value
}

func extract(logLine map[string]any, key string) (string, bool) {
	valueRaw, ok := logLine[key]
	if ok {
		delete(logLine, key)
	}
	value := asString(valueRaw)
	return value, ok
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
	delete(logLine, "level")

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
		f(a...)
	}
}
