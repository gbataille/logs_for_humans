package consolefmt

import (
	"sync"

	"github.com/gbataille/logs_for_humans/consolefmt/json"
	"github.com/gbataille/logs_for_humans/consolefmt/raw"
	"github.com/gbataille/logs_for_humans/in_stream"
	"github.com/gbataille/logs_for_humans/parsing"
	"github.com/pterm/pterm"
)

func handleLine(lineB []byte) {
	if len(lineB) == 0 {
		return
	}

	logLine, err := parsing.FromJsonGeneric(lineB)
	if err != nil {
		raw.HandleRAWLine(string(lineB))
	} else {
		json.HandleJSONLogLine(logLine)
	}
}

func handleError(err error) {
	pterm.ThemeDefault.ErrorMessageStyle.Println(err.Error())
}

func Run() {
	lineChan, errorChan := in_stream.HandleStdinByLine()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for line := range lineChan {
			handleLine(line)
		}
		wg.Done()
	}()

	go func() {
		for err := range errorChan {
			handleError(err)
		}
		wg.Done()
	}()

	wg.Wait()
}
