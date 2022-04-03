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

// HandleStdinLines reads continuously from STDIN and sends the lines read on the channel
// If it tries to parse a too long line, sends an error to the error handler and a non-finished line on the channel
func HandleStdinByLine() (chan []byte, chan error) {
	lineChan := make(chan []byte)
	errorChan := make(chan error)

	go func() {
		// It is critical to parse small chunks, because linux PIPE have a max buffer size
		// We quickly consume small batches and process them to extract individual lines
		buf := make([]byte, 4096)
		nextLine := make([]byte, 0, size65k)

		dumpAndReset := func() {
			lineChan <- nextLine
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
				errorChan <- errors.New("read buffer full without a newline, dumping")
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
			lineChan <- nextLine
		}

		close(lineChan)
		close(errorChan)
	}()

	return lineChan, errorChan
}
