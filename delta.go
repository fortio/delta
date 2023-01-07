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
	"sort"

	"fortio.org/fortio/log"
	"fortio.org/fortio/version"
)

var (
	fullVersion = flag.Bool("version", false, "Show full version info and exit.")
	aCmd        = flag.String("a", "", "`Command` to run for each entry unique to file A")
	bCmd        = flag.String("b", "", "`Command` to run for each entry unique to file B")
)

func usage(msg string) {
	_, _ = fmt.Fprintf(os.Stderr, "Fortio delta %s usage:\n\t%s [flags] fileA fileB\nflags:\n",
		version.Short(),
		os.Args[0])
	flag.PrintDefaults()
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

func main() {
	flag.CommandLine.Usage = func() { usage("") }
	log.SetFlagDefaultsForClientTools()
	flag.Parse()
	_, longV, fullV := version.FromBuildInfo()
	if *fullVersion {
		fmt.Print(fullV)
		os.Exit(0)
	}
	if len(flag.Args()) != 2 {
		usage("Need 2 arguments (old and new files)")
	}
	log.Infof("Fortio delta %s started - will run %q on new entries, and %q on missing ones", longV, *aCmd, *bCmd)
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
		if *aCmd != "" {
			log.Infof("Running %s %s", *aCmd, a)
		}
	}
	onlyInB := sortedKeys(bSet)
	for _, b := range onlyInB {
		log.Infof("Only in B: %q", b)
		if *bCmd != "" {
			log.Infof("Running %s %s", *bCmd, b)
		}
	}
}
