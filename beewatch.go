// Copyright 2013 Unknown
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Bee Watch is an interactive debugger for the Go programming language.
package beewatch

import (
	"encoding/json"
	"fmt"
	"os"
)

type debugLevel int

const (
	Trace debugLevel = iota
	Info
	Critical
)

var (
	watchLevel      debugLevel
	beewatchEnabled bool
)

var App struct {
	Name     string `json:"app_name"`
	HttpPort int    `json:"http_port"`
}

const (
	APP_VER = "0.1.0.0809"
)

// Start initialize debugger data.
func Start(wl debugLevel) {
	fmt.Printf("[INIT] BW: Bee Watch v%s.\n", APP_VER)
	loadJSON()
	watchLevel = wl
	beewatchEnabled = true
	initHTTP()
}

func loadJSON() {
	f, err := os.Open("beewatch.json")
	if err != nil {
		fmt.Printf("[ERRO] BW: Fail to load beewatch.json[ %s ]\n", err)
		os.Exit(2)
	}
	defer f.Close()

	d := json.NewDecoder(f)
	err = d.Decode(&App)
	if err != nil {
		fmt.Printf("[ERRO] BW: Fail to parse beewatch.json[ %s ]\n", err)
		os.Exit(2)
	}
}
