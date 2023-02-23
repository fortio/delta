package main

import (
	"net/url"
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"delta": Main,
	}))
}

func TestDelta(t *testing.T) {
	testscript.Run(t, testscript.Params{Dir: "./"})
}
