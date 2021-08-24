package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

// Details define a package's name and relevant security fixes included in a
// given version.
type Details struct {
	Name string `json:"name"`
	// Fixed package version string mapped to an array of CVE ids affecting the
	// package.
	Secfixes map[string][]string `json:"secfixes"`
}

// Package wraps the Details.
type SecdbPackage struct {
	Pkg Details `json:"pkg"`
}

// SecurityDB is the security database structure.
type SecurityDB struct {
	Distroversion string         `json:"distroversion"`
	Reponame      string         `json:"reponame"`
	Urlprefix     string         `json:"urlprefix"`
	Apkurl        string         `json:"apkurl"`
	Packages      []SecdbPackage `json:"packages"`
}

func fetchSecDb() (*SecurityDB, error) {
	client := http.Client{}
	resp, err := client.Get("https://secdb.alpinelinux.org/edge/community.json")
	if err != nil {
		return nil, fmt.Errorf("failed to request alpine secdb: %v", err)
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("failed to load alpine secdb: %v", err)
	}
	db := SecurityDB {}
	err = json.Unmarshal(body, &db)
	if err != nil {
		return nil, fmt.Errorf("failed to load alpine secdb: %v", err)
	}
	return &db, nil
}
