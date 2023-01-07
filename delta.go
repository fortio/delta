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
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"fortio.org/fortio/log"
	"fortio.org/fortio/version"
)

var (
	fullVersion = flag.Bool("version", false, "Show full version info and exit.")
	aCmd        = flag.String("a", "", "`Command` to run for each entry unique to file A")
	bCmd        = flag.String("b", "", "`Command` to run for each entry unique to file B")
	shortV      string
)

func usage(msg string) {
	_, _ = fmt.Fprintf(os.Stderr, "Fortio delta %s usage:\n\t%s [flags] fileA fileB\nflags:\n",
		shortV,
		os.Args[0])
	flag.PrintDefaults()
	_, _ = fmt.Fprintf(os.Stderr, "Command flags (-a and -b) are space separeted and the lines are passed as the last argument")
	if msg != "" {
		fmt.Fprintln(os.Stderr, msg)
	}
	os.Exit(1)
}

func toMap(filename string) (map[string]bool, error) {
	log.LogVf("Reading %q", filename)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m := make(map[string]bool)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		entry := sc.Text()
		log.LogVf("  adding %q", entry)
		m[entry] = true
	}
	return m, nil
}

func removeCommon(a, b map[string]bool) {
	if len(a) > len(b) {
		a, b = b, a
	}
	for e := range a {
		if _, found := b[e]; found {
			log.LogVf("in both sets: %q", e)
			delete(a, e)
			delete(b, e)
		}
	}
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func runCmd(cmd0 string, args ...string) {
	cmd := exec.Command(cmd0, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Infof("Running (%d args) %s", len(args), cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error running %v: %v", cmd, err)
	}
}

func main() {
	flag.CommandLine.Usage = func() { usage("") }
	log.SetFlagDefaultsForClientTools()
	sV, longV, fullV := version.FromBuildInfo()
	shortV = sV
	flag.Parse()
	if *fullVersion {
		fmt.Print(fullV)
		os.Exit(0)
	}
	if len(flag.Args()) != 2 {
		usage("Need 2 arguments (old and new files)")
	}
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
		longV, *aCmd, flag.Arg(0), *bCmd, flag.Arg(1))
	// read file content into map
	aSet, err := toMap(flag.Arg(0))
	if err != nil {
		log.Fatalf("Error reading old file: %v", err)
	}
	bSet, err := toMap(flag.Arg(1))
	if err != nil {
		log.Fatalf("Error reading new file: %v", err)
	}
	// remove common entries
	removeCommon(aSet, bSet)
	onlyInA := sortedKeys(aSet)
	for _, a := range onlyInA {
		log.Infof("Only in A: %q", a)
		if hasACmd {
			aArgs[aLen] = a
			runCmd(aCmd0, aArgs...)
		}
	}
	onlyInB := sortedKeys(bSet)
	for _, b := range onlyInB {
		log.Infof("Only in B: %q", b)
		if hasBCmd {
			bArgs[bLen] = b
			runCmd(bCmd0, bArgs...)
		}
	}
}
