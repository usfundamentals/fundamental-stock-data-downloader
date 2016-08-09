// Copyright 2016 Linelane GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
)

const usageText = `Usage example:

usf-downloader update -auth-token=YOUR_TOKEN -data-dir=./data

`

func help() {
	fmt.Println(usageText)
	flag.PrintDefaults()
}

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var authToken = flags.String("auth-token", "", "Auth token. Get one at https://account.usfundamentals.com")
	var dataDir = flags.String("data-dir", "", "Dir for storing downloaded data.")

	flag.Usage = help

	if len(os.Args) < 2 {
		help()
		os.Exit(2)
	}

	flags.Parse(os.Args[2:])

	action := os.Args[1]

	switch action {
	case "update":
		if *authToken == "" || *dataDir == "" {
			help()
			os.Exit(2)
		}
		actionUpdate(*dataDir, *authToken)
	default:
		help()
		os.Exit(2)
	}
}
