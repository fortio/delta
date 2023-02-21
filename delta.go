// Copyright 2023 Fortio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"flag"
	"os"
	"os/exec"
	"strings"

	"fortio.org/cli"
	"fortio.org/log"
	"fortio.org/sets"
)

var (
	aCmd = flag.String("a", "", "`Command` to run for each entry unique to file A")
	bCmd = flag.String("b", "", "`Command` to run for each entry unique to file B")
)

func toMap(filename string) (sets.Set[string], error) {
	log.LogVf("Reading %q", filename)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m := sets.New[string]()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		entry := sc.Text()
		log.LogVf("  adding %q", entry)
		m.Add(entry)
	}
	return m, nil
}

func runCmd(cmd0 string, args ...string) bool {
	cmd := exec.Command(cmd0, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Infof("Running (%d args) %s", len(args), cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Errf("Error running %v: %v", cmd, err)
		return false
	}
	return true
}

func main() {
	os.Exit(Main())
}

// Main is the main function for the delta tool so it can be called from testscript.
// Note that we could use the (new in 1.39) log.Fatalf that doesn't panic for cli tools but
// it calling os.Exit directly means it doesn't work with the code coverage from `testscript`
// but there is now (in 1.40) log.FErrf for that (no exit Fatalf).
func Main() int {
	cli.ProgramName = "Fortio delta"
	cli.MinArgs = 2
	cli.ArgsHelp = "fileA fileB" +
		"\nwith command flags (-a and -b) space separated and the lines are passed as the last argument"
	cli.Main()
	cmdList := strings.Split(*aCmd, " ")
	aCmd0 := cmdList[0]
	aArgs := cmdList[1:]
	aLen := len(aArgs)
	aArgs = append(aArgs, "ONLY_IN_A") // placeholder for later
	hasACmd := len(aCmd0) > 0
	cmdList = strings.Split(*bCmd, " ")
	bCmd0 := cmdList[0]
	bArgs := cmdList[1:]
	bLen := len(bArgs)
	bArgs = append(bArgs, "ONLY_IN_A") // placeholder for later
	hasBCmd := len(bCmd0) > 0
	log.Infof("Fortio delta %s started - will run %q on entries unique to %s, and %q on ones unique to %s",
		cli.LongVersion, *aCmd, flag.Arg(0), *bCmd, flag.Arg(1))
	// read file content into map
	aSet, err := toMap(flag.Arg(0))
	if err != nil {
		return log.FErrf("Error reading file A: %v", err)
	}
	bSet, err := toMap(flag.Arg(1))
	if err != nil {
		return log.FErrf("Error reading file B: %v", err)
	}
	// remove common entries
	sets.RemoveCommon(aSet, bSet)
	onlyInA := aSet.Sorted()
	for _, a := range onlyInA {
		log.Infof("Only in A: %q", a)
		if hasACmd {
			aArgs[aLen] = a
			if !runCmd(aCmd0, aArgs...) {
				return 1
			}
		}
	}
	onlyInB := bSet.Sorted()
	for _, b := range onlyInB {
		log.Infof("Only in B: %q", b)
		if hasBCmd {
			bArgs[bLen] = b
			if !runCmd(bCmd0, bArgs...) {
				return 1
			}
		}
	}
	return 0
}
