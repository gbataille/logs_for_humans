package json

import (
	"fmt"
	"sort"
	"strings"

	pterm_utils "github.com/gbataille/logs_for_humans/utils/pterm"
	"github.com/pterm/pterm"
)

const maxPanelWidth int = 50

func HandleJSONLogLine(logLine map[string]interface{}) {
	logFunc := methodFromLine(logLine)
	logMsg, ok := extract(logLine, "msg")

	if !ok {
		logFunc(fmt.Sprintf("%v", logLine))
		return
	}

	width, _, err := pterm.GetTerminalSize()
	if err != nil {
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

	var maxKeySize int
	keys := make([]string, 0, len(logLine))
	for k := range logLine {
		keys = append(keys, k)
		if len(k) > maxKeySize {
			maxKeySize = len(k)
		}
	}
	sort.Strings(keys)

	// tables uses 3 chars as a separator
	// Plus we indent it with 4 spaces at the start
	// We have room for nbTables tables, with a width of maxPanelWidth
	nbTables := int((width - 4) / (maxPanelWidth + 3))
	actualTableWidth := int((width - 4) / nbTables)

	maxValueSpace := actualTableWidth - 3 - maxKeySize

	tables := make([][][]string, 0, nbTables)
	for i := 0; i < nbTables; i++ {
		tables = append(tables, make([][]string, 0, len(logLine)/nbTables))
	}

	rawOut := make(map[string]string)

	for idx, k := range keys {
		v := asString(logLine[k])
		if len(v) > maxValueSpace {
			rawOut[k] = v
		} else {
			item := []string{
				pterm.ThemeDefault.DebugMessageStyle.Sprintf("    %s", k),
				pterm.ThemeDefault.DebugMessageStyle.Sprint(asString(v)),
			}
			tables[idx%nbTables] = append(tables[idx%nbTables], item)
		}
	}

	panels := make([]pterm.Panel, 0, nbTables)
	for _, table := range tables {
		data, _ := pterm.DefaultTable.WithData(pterm.TableData(table)).Srender()
		panels = append(panels, pterm.Panel{Data: data})
	}
	ps := pterm.Panels{panels}
	err = pterm.DefaultPanel.WithPanels(ps).WithPadding(5).Render()
	if err != nil {
		fmt.Println("ERROR ########")
	}

	for k, v := range rawOut {
		k = "    " + k + " :"
		pterm.ThemeDefault.DebugMessageStyle.Println(k)
		v = strings.ReplaceAll(v, "\n", "\n        ")
		v = "        " + v
		pterm.ThemeDefault.DebugMessageStyle.Println(v)
	}
}

func asString(raw interface{}) string {
	var value string
	switch raw := raw.(type) {
	case string:
		value = raw
	case fmt.Stringer:
		value = raw.String()
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

	return pterm_utils.MethodFromLevel(level)
}
