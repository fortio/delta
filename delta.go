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

	"fortio.org/fortio/log"
	"fortio.org/fortio/version"
)

var (
	fullVersion = flag.Bool("version", false, "Show full version info and exit.")
	newCmd      = flag.String("new", "", "`Command` to run for each entry unique to new file")
	goneCmd     = flag.String("gone", "", "`Command` to run for each entry missing in new file")
)

func usage(msg string) {
	_, _ = fmt.Fprintf(os.Stderr, "Fortio delta %s usage:\n\t%s [flags] old new\nflags:\n",
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

func main() {
	flag.CommandLine.Usage = func() { usage("") }
	flag.Parse()
	_, longV, fullV := version.FromBuildInfo()
	if *fullVersion {
		fmt.Print(fullV)
		os.Exit(0)
	}
	if len(flag.Args()) != 2 {
		usage("Need 2 arguments (old and new files)")
	}
	log.Infof("Fortio delta %s started - will run %q on new entries, and %q on missing ones", longV, *newCmd, *goneCmd)
	// read file content into map
	oldSet, err := toMap(flag.Arg(0))
	if err != nil {
		log.Fatalf("Error reading old file: %v", err)
	}
	newSet, err := toMap(flag.Arg(1))
	if err != nil {
		log.Fatalf("Error reading new file: %v", err)
	}
	// intersect
	for o := range oldSet {
		if _, found := newSet[o]; found {
			log.LogVf("old loop: %q is in both", o)
			// mark as present in both - tradeoff of mutating to avoid extra lookup later (even though o(1))
			// check if mutating the value is better than removing the entry (todo benchmark)
			newSet[o] = false
			continue
		}
		log.Infof("Missing %q", o)
		if *goneCmd != "" {
			log.Infof("Running %s %s", *goneCmd, o)
		}
	}
	// is there a more efficient way to get only the "true" keys
	for n, v := range newSet {
		if !v {
			log.LogVf("new loop: %q already known to be in both", n)
			continue
		}
		log.Infof("New %q", n)
		if *newCmd != "" {
			log.Infof("Running %s %s", *newCmd, n)
		}
	}
}
