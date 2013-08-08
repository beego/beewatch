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
	"os"
	"runtime"
	"strings"
)

const (
	IMPORT_PATH = "github.com/beego/beewatch"
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
		sp := gp + "/src/" + IMPORT_PATH + "/static"
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
	case Trace:
		return "TRACE"
	case Info:
		return "INFO"
	case Critical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}
