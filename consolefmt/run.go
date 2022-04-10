package consolefmt

import (
	"sync"

	"github.com/gbataille/logs_for_humans/consolefmt/json"
	"github.com/gbataille/logs_for_humans/consolefmt/raw"
	"github.com/gbataille/logs_for_humans/in_stream"
	"github.com/gbataille/logs_for_humans/parsing"
	"github.com/pterm/pterm"
)

type lineType int

const (
	unknown     lineType = -1
	rawLineType lineType = iota
	jsonLineType
)

type ConsoleFormatter struct {
	lastLine chan lineType
}

func (cf *ConsoleFormatter) readLastLineTypeNonBlocking() lineType {
	select {
	case lastLineType := <-cf.lastLine:
		return lastLineType
	default:
		return -1
	}
}

func (cf *ConsoleFormatter) handleLine(lineB []byte) {
	if len(lineB) == 0 {
		return
	}

	// Always consume whatever has been put in the channel before.
	// Does not block because the channel starts empty
	lastLineType := cf.readLastLineTypeNonBlocking()

	logLine, err := parsing.FromJsonGeneric(lineB)
	if err != nil {
		cf.lastLine <- rawLineType
		raw.HandleRAWLine(string(lineB))
	} else {
		// Leave a blank line after a group of raw lines
		if lastLineType == rawLineType {
			pterm.Println()
		}

		cf.lastLine <- jsonLineType
		json.HandleJSONLogLine(logLine)
	}
}

func (cf *ConsoleFormatter) handleError(err error) {
	pterm.ThemeDefault.ErrorMessageStyle.Println(err.Error())
}

func (cf *ConsoleFormatter) Run() {
	cf.lastLine = make(chan lineType, 1)

	lineChan, errorChan := in_stream.HandleStdinByLine()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for line := range lineChan {
			cf.handleLine(line)
		}
		wg.Done()
	}()

	go func() {
		for err := range errorChan {
			cf.handleError(err)
		}
		wg.Done()
	}()

	wg.Wait()
}
