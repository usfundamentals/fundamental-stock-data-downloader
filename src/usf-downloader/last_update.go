// Copyright 2016 Linelane GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const lastFullDownloadTimeFile = "last-full-download-time"
const lastUpdateIDFile = "last-update-id"

func setLastFullDownloadTime(dataDir string, ts time.Time) {
	data := []byte(strconv.FormatInt(ts.UnixNano(), 10))
	err := ioutil.WriteFile(path.Join(dataDir, lastFullDownloadTimeFile), data, 0777)
	if err != nil {
		panic(err)
	}
}

func getLastFullDownloadTime(dataDir string) time.Time {
	data, err := ioutil.ReadFile(path.Join(dataDir, lastFullDownloadTimeFile))
	if os.IsNotExist(err) {
		return time.Time{}
	}
	if err != nil {
		panic(err)
	}
	ns, err := strconv.ParseInt(strings.Trim(string(data), "\n"), 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(0, ns)
}

func setLastUpdateID(dataDir string, id string) {
	err := ioutil.WriteFile(path.Join(dataDir, lastUpdateIDFile), []byte(id), 0777)
	if err != nil {
		panic(err)
	}
}

func getLastUpdateID(dataDir string) string {
	data, err := ioutil.ReadFile(path.Join(dataDir, lastUpdateIDFile))
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(data), "\n")
}
