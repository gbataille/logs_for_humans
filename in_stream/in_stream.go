package in_stream

import (
	"errors"
	"io"
	"math"
	"os"
)

const newLine byte = 10

var size65k int
var sizeMax int

func init() {
	size65k = int(math.Pow(2, 16))
	sizeMax = int(math.Pow(2, 20))
}

type StdinLineHandler interface {
	HandleLine(line []byte)
	HandleError(err error)
}

// HandleStdinLines reads continuously from STDIN and calls the function passed in parameter for each line of data received
// If it tries to parse a too long line, sends an error to the error handler
func HandleStdinByLine(handler StdinLineHandler) {

	// It is critical to parse small chunks, because linux PIPE have a max buffer size
	// We quickly consume small batches and process them to extract individual lines
	buf := make([]byte, 4096)
	nextLine := make([]byte, 0, size65k)

	dumpAndReset := func() {
		handler.HandleLine(nextLine)
		nextLine = make([]byte, 0, 2^16)
	}

	n, err := os.Stdin.Read(buf)
	// Loop until STDIN closes
	for err == nil {
		for _, char := range buf[:n] {
			if char == newLine {
				dumpAndReset()
			} else {
				nextLine = append(nextLine, char)
			}
		}

		if len(nextLine) >= sizeMax {
			handler.HandleError(errors.New("read buffer full without a newline, dumping"))
			dumpAndReset()
		}

		// Loop increment
		n, err = os.Stdin.Read(buf)
	}

	if err != io.EOF {
		// As far as I know, this is not "possible"
		panic(err.Error())
	} else {
		// Handle the last line consumed
		handler.HandleLine(nextLine)
	}
}
