package main

import (
	"os"
	"testing"

	"github.com/test-go/testify/require"
)

func TestFailingLine(t *testing.T) {
	data, err := os.ReadFile("./samples/failing_line.json")
	require.NoError(t, err)
	handleLine(data)
}
