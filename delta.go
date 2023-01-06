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
	"flag"
	"fmt"
	"os"

	"fortio.org/fortio/log"
	"fortio.org/fortio/version"
)

var fullVersion = flag.Bool("version", false, "Show full version info and exit.")

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
	log.Infof("Fortio delta %s started", longV)
}
