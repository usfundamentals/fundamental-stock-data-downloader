// Copyright 2016 Linelane GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"os"
	"path"
)

func applyUpdate(dataDir, companyID, period string, updateData []byte) {
	data, err := readIndicators(dataDir, companyID)
	if os.IsNotExist(err) {
		writeIndicators(dataDir, companyID, updateData)
		return
	}

	current := parseIndicators(bytesToRows(data))
	update := parseIndicators(bytesToRows(updateData))

	// update data may be empty in this case we need to delete the data for period

	// delete current values
	for _, d := range current {
		delete(d, period)
	}

	// update with new values
	for indicator, d := range update {
		for period, value := range d {
			if _, ok := current[indicator]; !ok {
				current[indicator] = map[string]string{}
			}
			current[indicator][period] = value
		}
	}

	rows := indicatorMapToRows(current, companyID)
	csvData := rowsToBytes(rows)
	writeIndicators(dataDir, companyID, csvData)
}

func parseIndicators(rows [][]string) indicatorMap {
	res := indicatorMap{}
	if len(rows) <= 1 { // no data rows
		return res
	}
	if len(rows[0]) <= 2 { // no periods
		return res
	}
	periods := rows[0][2:]
	for _, row := range rows[1:] {
		indicator := row[1]
		values := map[string]string{}
		for i, v := range row[2:] {
			values[periods[i]] = v
		}
		res[indicator] = values
	}
	return res
}

// indicatorMap map[indictor]map[period]value
type indicatorMap map[string]map[string]string

func bytesToRows(data []byte) [][]string {
	r := bytes.NewReader(data)
	csvr := csv.NewReader(r)
	rows, err := csvr.ReadAll()
	if err != nil {
		panic(err)
	}
	return rows
}

func rowsToBytes(rows [][]string) []byte {
	wr := &bytes.Buffer{}
	csvw := csv.NewWriter(wr)
	err := csvw.WriteAll(rows)
	if err != nil {
		panic(err)
	}
	return wr.Bytes()
}

func readIndicators(dataDir, companyID string) ([]byte, error) {
	indicatorsDir := path.Join(dataDir, "indicators_by_company")
	return ioutil.ReadFile(path.Join(indicatorsDir, companyID+".csv"))
}

func writeIndicators(dataDir, companyID string, data []byte) {
	indicatorsDir := path.Join(dataDir, "indicators_by_company")
	os.Mkdir(indicatorsDir, 0777)

	err := ioutil.WriteFile(path.Join(indicatorsDir, companyID+".csv"), data, 0666)
	if err != nil {
		panic(err)
	}
}
