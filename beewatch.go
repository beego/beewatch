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
	"os"
)

const (
	APP_VER = "0.4.5.0811"
)

type debugLevel int

const (
	LevelTrace debugLevel = iota
	LevelInfo
	LevelCritical
)

var (
	watchLevel debugLevel
)

var App struct {
	Name         string `json:"app_name"`
	HttpPort     int    `json:"http_port"`
	WatchEnabled bool   `json:"watch_enabled"`
	CmdMode      bool   `json:"cmd_mode"`
}

// Start initialize debugger data.
func Start(wl ...debugLevel) {
	colorLog("[INIT] BW: Bee Watch v%s.\n", APP_VER)

	watchLevel = LevelTrace
	if len(wl) > 0 {
		watchLevel = wl[0]
	}

	App.Name = "Bee Watch"
	App.HttpPort = 23456
	App.WatchEnabled = true

	loadJSON()

	if App.WatchEnabled && !App.CmdMode {
		initHTTP()
	}
}

func loadJSON() {
	f, err := os.Open("beewatch.json")
	if err != nil {
		colorLog("[ERRO] BW: Fail to load beewatch.json[ %s ]\n", err)
		os.Exit(2)
	}
	defer f.Close()

	d := json.NewDecoder(f)
	err = d.Decode(&App)
	if err != nil {
		colorLog("[ERRO] BW: Fail to parse beewatch.json[ %s ]\n", err)
		os.Exit(2)
	}
}
