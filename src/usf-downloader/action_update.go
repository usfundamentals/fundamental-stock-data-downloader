// Copyright 2016 Linelane GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func actionUpdate(dataDir, token string) {
	_, err := os.Stat(dataDir)
	if os.IsNotExist(err) {
		fmt.Println("Specified data dir does not exist, please create it first.")
		os.Exit(1)
	}

	if getLastFullDownloadTime(dataDir).IsZero() {
		initialDownload(dataDir, token)
		return
	}

	updateData(dataDir, token)
}

const apiRoot = "https://api.usfundamentals.com"
const apiCompanies = apiRoot + "/v1/companies/xbrl"

type companyJSON struct {
	CompanyID     string   `json:"company_id"`
	NameLatest    string   `json:"name_latest"`
	NamesPrevious []string `json:"names_previous"`
}

func initialDownload(dataDir, token string) {

	fmt.Println("starting initial download")

	startTime := time.Now()

	companies := getCompanies(token)
	for i, item := range companies {
		fmt.Println("downloading data for company", item.CompanyID, i, "of", len(companies))
		data := getIndicators(token, item.CompanyID, "")
		writeIndicators(dataDir, item.CompanyID, data)
	}

	setLastFullDownloadTime(dataDir, startTime)
}

func getCompanies(token string) []companyJSON {
	resp, err := http.Get(apiCompanies + "?format=json&token=" + token)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var companies []companyJSON
	err = json.Unmarshal(data, &companies)
	if err != nil {
		panic(err)
	}
	return companies
}

const apiIndicators = apiRoot + "/v1/indicators/xbrl"

func getIndicators(token string, company, period string) []byte {
	query := "?token=" + token + "&companies=" + company

	resp, err := http.Get(apiIndicators + query)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return data
}

func getIndicatorsPeriod(token string, company, period string) []byte {
	query := "?token=" + token + "&companies=" + company + "&periods=" + period

	resp, err := http.Get(apiIndicators + query)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return data
}

const apiUpdates = apiRoot + "/v1/indicators/xbrl/updates"

type updateJSON struct {
	UpdateID  string `json:"update_id"`
	CompanyID string `json:"company_id"`
	Period    string `json:"period"`
}

func getUpdates(token string, fromTime time.Time, fromID string) []updateJSON {
	query := "?token=" + token
	if fromID != "" {
		query += "&update_id_from=" + fromID
	} else if !fromTime.IsZero() {
		query += "&updated_on_from=" + strconv.FormatInt(fromTime.UnixNano(), 10)
	} else {
		panic("provide fromID or fromTime")
	}

	resp, err := http.Get(apiUpdates + query)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var updates []updateJSON
	err = json.Unmarshal(data, &updates)
	if err != nil {
		panic(err)
	}
	return updates
}

func updateData(dataDir, token string) {
	fmt.Println("starting update")

	var updates []updateJSON

	id := getLastUpdateID(dataDir)
	downloadTime := getLastFullDownloadTime(dataDir)
	if id != "" {
		updates = getUpdates(token, time.Time{}, id)
	} else if !downloadTime.IsZero() {
		updates = getUpdates(token, downloadTime, "")
	} else {
		panic("no existing data found, run initialDownload first")
	}

	if len(updates) == 0 {
		fmt.Println("data is up-to-date")
		return
	}

	for _, update := range updates {
		fmt.Println("updating company", update.CompanyID, "period", update.Period)
		updateData := getIndicatorsPeriod(token, update.CompanyID, update.Period)
		applyUpdate(dataDir, update.CompanyID, update.Period, updateData)
		setLastUpdateID(dataDir, update.UpdateID)
	}
}
