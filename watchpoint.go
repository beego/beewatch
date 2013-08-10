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
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
)

var debuggerMutex = sync.Mutex{}

type WatchPoint struct {
	disabled   bool
	watchLevel debugLevel
	offset     int
}

// Trace debug level.
func Trace() *WatchPoint {
	return setlevel(LevelTrace)
}

// Info debug level.
func Info() *WatchPoint {
	return setlevel(LevelInfo)
}

// Critical debug level.
func Critical() *WatchPoint {
	return setlevel(LevelCritical)
}

func setlevel(wl debugLevel) *WatchPoint {
	disabled := false
	if !beewatchEnabled || watchLevel >= LevelInfo {
		disabled = true
	}

	return &WatchPoint{
		disabled:   disabled,
		watchLevel: wl,
		offset:     1,
	}
}

// Display sends variable name,value pairs to the debugger. Values are formatted using %#v.
// The parameter 'nameValuePairs' must be even sized.
func (wp *WatchPoint) Display(nameValuePairs ...interface{}) *WatchPoint {
	if wp.disabled {
		return wp
	}

	_, file, line, ok := runtime.Caller(wp.offset)
	cmd := command{
		Action: "DISPLAY",
		Level:  levelToStr(wp.watchLevel),
	}
	if ok {
		cmd.addParam("go.file", file)
		cmd.addParam("go.line", fmt.Sprint(line))
	}
	if len(nameValuePairs)%2 == 0 {
		for i := 0; i < len(nameValuePairs); i += 2 {
			k := nameValuePairs[i]
			v := nameValuePairs[i+1]
			cmd.addParam(fmt.Sprint(k), fmt.Sprintf("%#v", v))
		}
	} else {
		fmt.Printf("[WARN] BW: Missing variable for Display(...) in: %v:%v.\n", file, line)
		wp.disabled = true
		return wp
	}

	channelExchangeCommands(wp.watchLevel, cmd)
	return wp
}

// Break halts the execution of the program and waits for an instruction from the debugger (e.g. Resume).
// Break is only effective if all (if any) conditions are true. The program will resume otherwise.
func (wp *WatchPoint) Break(conditions ...bool) *WatchPoint {
	if wp.disabled {
		return wp
	}

	suspend(wp, conditions...)
	return wp
}

// suspend will create a new Command and send it to the browser.
func suspend(wp *WatchPoint, conditions ...bool) {
	for _, condition := range conditions {
		if !condition {
			return
		}
	}

	_, file, line, ok := runtime.Caller(wp.offset)
	cmd := command{
		Action: "BREAK",
		Level:  levelToStr(wp.watchLevel),
	}
	if ok {
		cmd.addParam("go.file", file)
		cmd.addParam("go.line", fmt.Sprint(line))
		cmd.addParam("go.stack", trimStack(string(debug.Stack())))
	}
	channelExchangeCommands(wp.watchLevel, cmd)
}

// Peel off the part of the stack that lives in hopwatch
func trimStack(stack string) string {
	lines := strings.Split(stack, "\n")
	c := 0
	for _, line := range lines {
		// means no function in this package.
		if strings.Index(line, "/beewatch") == -1 {
			break
		}
		c++
	}
	return strings.Join(lines[c:], "\n")
}

func (wp *WatchPoint) Printf(format string, params ...interface{}) *WatchPoint {
	wp.offset += 1
	var content string
	if len(params) == 0 {
		content = format
	} else {
		content = fmt.Sprintf(format, params...)
	}
	return wp.printcontent(content)
}

// Printf formats according to a format specifier and writes to the debugger screen.
func (wp *WatchPoint) printcontent(content string) *WatchPoint {
	_, file, line, ok := runtime.Caller(wp.offset)
	cmd := command{
		Action: "PRINT",
		Level:  levelToStr(wp.watchLevel),
	}
	if ok {
		cmd.addParam("go.file", file)
		cmd.addParam("go.line", fmt.Sprint(line))
	}
	cmd.addParam("PRINT", content)
	channelExchangeCommands(wp.watchLevel, cmd)
	return wp
}

// Put a command on the browser channel and wait for the reply command.
func channelExchangeCommands(wl debugLevel, toCmd command) {
	if !beewatchEnabled || wl < watchLevel {
		return
	}

	// synchronize command exchange ; break only one goroutine at a time.
	debuggerMutex.Lock()
	toBrowserChannel <- toCmd
	<-fromBrowserChannel
	debuggerMutex.Unlock()
}
