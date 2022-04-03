package consolefmt

import (
	"github.com/gbataille/logs_for_humans/consolefmt/json"
	"github.com/gbataille/logs_for_humans/consolefmt/raw"
	"github.com/gbataille/logs_for_humans/in_stream"
	"github.com/gbataille/logs_for_humans/parsing"
	"github.com/pterm/pterm"
)

type Handler struct{}

func (h *Handler) HandleLine(lineB []byte) {
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

func (h *Handler) HandleError(err error) {
	pterm.ThemeDefault.ErrorMessageStyle.Println(err.Error())
}

func Run() {
	handler := &Handler{}
	in_stream.HandleStdinByLine(handler)
}
