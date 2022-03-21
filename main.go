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

func fromJsonGeneric(in []byte) (map[string]interface{}, error) {
	jsonMap := make(map[string](interface{}))
	err := json.Unmarshal(in, &jsonMap)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}

func toHumanLog(logLine map[string]interface{}) {
	logFunc := methodFromLine(logLine)
	logMsg, ok := extract(logLine, "msg")

	if !ok {
		logFunc(fmt.Sprintf("%v", logLine))
		return
	}

	logTs, _ := extract(logLine, "ts")
	logFunc(fmt.Sprintf("[%s]", logTs), logMsg)

	caller, _ := extract(logLine, "caller")
	caller = pterm.ThemeDefault.DescriptionMessageStyle.Sprint(caller)
	leveledList := pterm.LeveledList{
		pterm.LeveledListItem{Level: 0, Text: caller},
	}
	root := pterm.NewTreeFromLeveledList(leveledList)
	tree, _ := pterm.DefaultTree.WithRoot(root).Srender()
	pterm.Print(tree)

	rawKeys := []string{"panic_stack_trace"}
	rawOut := make([]string, 0, len(rawKeys))
	for _, key := range rawKeys {
		out, ok := extract(logLine, key)
		if ok {
			rawOut = append(rawOut, out)
		}
	}

	keys := make([]string, 0, len(logLine))
	for k := range logLine {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	nbPanel := 3
	tables := make([][][]string, 0, nbPanel)
	for i := 0; i < nbPanel; i++ {
		tables = append(tables, make([][]string, 0, len(logLine)/nbPanel))
	}

	for idx, k := range keys {
		item := []string{
			pterm.ThemeDefault.DebugMessageStyle.Sprintf("\t\t%s", k),
			pterm.ThemeDefault.DebugMessageStyle.Sprint(asString(logLine[k])),
		}
		tables[idx%nbPanel] = append(tables[idx%nbPanel], item)
	}

	panels := make([]pterm.Panel, 0, nbPanel)
	for _, table := range tables {
		data, _ := pterm.DefaultTable.WithData(pterm.TableData(table)).Srender()
		panels = append(panels, pterm.Panel{Data: data})
	}
	ps := pterm.Panels{panels}
	pterm.DefaultPanel.WithPanels(ps).WithPadding(5).Render()

	for _, out := range rawOut {
		out = strings.ReplaceAll(out, "\n", "\n        ")
		out = "        " + out
		pterm.ThemeDefault.DebugMessageStyle.Println(out)
	}
}

func asString(raw interface{}) string {
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

func extract(logLine map[string]interface{}, key string) (string, bool) {
	valueRaw, ok := logLine[key]
	if ok {
		delete(logLine, key)
	}
	value := asString(valueRaw)
	return value, ok
}

func methodFromLine(logLine map[string]interface{}) func(a ...interface{}) {
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

func methodFromLevel(level string) func(a ...interface{}) {
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

func withNoReturn(f func(a ...interface{}) *pterm.TextPrinter) func(a ...interface{}) {
	return func(a ...interface{}) {
		f(a...)
	}
}
