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

package beewatch

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	_IMPORT_PATH = "github.com/beego/beewatch"
)

// isExist returns if a file or directory exists
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// getStaticPath returns app static path in somewhere in $GOPATH,
// it returns error if it cannot find.
func getStaticPath() (string, error) {
	gps := os.Getenv("GOPATH")

	var sep string
	if runtime.GOOS == "windows" {
		sep = ";"
	} else {
		sep = ":"
	}

	for _, gp := range strings.Split(gps, sep) {
		sp := gp + "/src/" + _IMPORT_PATH + "/static"
		if isExist(sp) {
			return sp, nil
		}
	}

	return "", errors.New("Cannot find static path in $GOPATH")
}

// loadFile loads file data by given path.
func loadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	b := make([]byte, fi.Size())
	f.Read(b)
	return b, nil
}

// levelToStr returns string format of debug level.
func levelToStr(wl debugLevel) string {
	switch wl {
	case LevelTrace:
		return "TRACE"
	case LevelInfo:
		return "INFO"
	case LevelCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

const (
	Gray = uint8(iota + 90)
	Red
	Green
	Yellow
	Blue
	Magenta
	//NRed      = uint8(31) // Normal
	EndColor = "\033[0m"
)

// colorLog colors log and print to stdout.
// Log format: <level> <content [highlight][path]> [ error ].
// Level: TRAC -> blue; ERRO -> red; WARN -> Magenta; SUCC -> green; others -> default.
// Content: default; path: yellow; error -> red.
// Level has to be surrounded by "[" and "]".
// Highlights have to be surrounded by "# " and " #"(space).
// Paths have to be surrounded by "( " and " )"(sapce).
// Errors have to be surrounded by "[ " and " ]"(space).
func colorLog(format string, a ...interface{}) {
	log := fmt.Sprintf(format, a...)
	if runtime.GOOS != "windows" {
		var clog string

		// Level.
		i := strings.Index(log, "]")
		if log[0] == '[' && i > -1 {
			clog += "[" + getColorLevel(log[1:i]) + "]"
		}

		log = log[i+1:]

		// Error.
		log = strings.Replace(log, "[ ", fmt.Sprintf("[\033[%dm", Red), -1)
		log = strings.Replace(log, " ]", EndColor+"]", -1)

		// Path.
		log = strings.Replace(log, "( ", fmt.Sprintf("(\033[%dm", Yellow), -1)
		log = strings.Replace(log, " )", EndColor+")", -1)

		// Highlights.
		log = strings.Replace(log, "# ", fmt.Sprintf("\033[%dm", Gray), -1)
		log = strings.Replace(log, " #", EndColor, -1)

		log = clog + log
	}

	fmt.Print(log)
}

// getColorLevel returns colored level string by given level.
func getColorLevel(level string) string {
	level = strings.ToUpper(level)
	switch level {
	case "TRAC":
		return fmt.Sprintf("\033[%dm%s\033[0m", Blue, level)
	case "ERRO":
		return fmt.Sprintf("\033[%dm%s\033[0m", Red, level)
	case "WARN":
		return fmt.Sprintf("\033[%dm%s\033[0m", Magenta, level)
	case "SUCC":
		return fmt.Sprintf("\033[%dm%s\033[0m", Green, level)
	default:
		return level
	}
}
