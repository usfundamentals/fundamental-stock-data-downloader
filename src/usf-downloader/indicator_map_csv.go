// Copyright 2016 Linelane GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"sort"

	"github.com/bradfitz/slice"
)

func indicatorMapToRows(m indicatorMap, companyID string) [][]string {

	periodsAll := map[string]bool{}
	for _, d := range m {
		for period := range d {
			periodsAll[period] = true
		}
	}
	periods := keysSB(periodsAll)
	sort.Strings(periods)

	rows := [][]string{}
	for indicator, d := range m {
		row := []string{companyID, indicator}
		values := []string{}
		for period := range periodsAll {
			values = append(values, d[period])
		}
		row = append(row, values...)
		rows = append(rows, row)
	}

	slice.Sort(rows, func(i, j int) bool {
		if rows[i][0] == rows[j][0] {
			return rows[i][1] < rows[j][1]
		}
		return rows[i][0] < rows[j][0]
	})

	headers := []string{"company_id", "indicator_id"}
	headers = append(headers, periods...)

	res := [][]string{headers}
	res = append(res, rows...)
	return res

}

func keysSB(m map[string]bool) []string {
	res := make([]string, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}
