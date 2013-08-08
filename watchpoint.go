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
	"sync"
)

var debuggerMutex = sync.Mutex{}

type WatchPoint struct {
	disabled   bool
	watchLevel debugLevel
	offset     int
}

// Display sends variable name:value pairs to the debugger.
// The parameter 'nameValuePairs' must be even sized.
func Display(wl debugLevel, nameValuePairs ...interface{}) *WatchPoint {
	if !beewatchEnabled || wl < watchLevel {
		return &WatchPoint{disabled: true}
	}

	wp := &WatchPoint{
		watchLevel: wl,
		offset:     2,
	}
	return wp.Display(wl, nameValuePairs...)
}

// Display sends variable name,value pairs to the debugger. Values are formatted using %#v.
// The parameter 'nameValuePairs' must be even sized.
func (wp *WatchPoint) Display(wl debugLevel, nameValuePairs ...interface{}) *WatchPoint {
	if wp.disabled {
		return wp
	}

	_, file, line, ok := runtime.Caller(wp.offset)
	cmd := command{
		Action: "DISPLAY",
		Level:  levelToStr(wl),
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
	channelExchangeCommands(wl, cmd)
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
