package main

import (
	"os"
	"testing"

	"fortio.org/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"delta": Main,
	}))
}

func TestDelta(t *testing.T) {
	testscript.Run(t, testscript.Params{Dir: "./"})
}
