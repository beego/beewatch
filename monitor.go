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
	//"bytes"
	"fmt"
	"reflect"
	"sync"
)

type watchValue struct {
	Value reflect.Value
	Note  string
}

var (
	watchLock *sync.RWMutex
	watchMap  map[string]*watchValue
)

// AddWatchVars adds variables to be watched,
// so that when the program calls 'Break()'
// these variables' infomation will be showed on monitor panel.
// The parameter 'nameValuePairs' must be even sized.
// Note that you have to pass pointers in order to have correct results.
func AddWatchVars(nameValuePairs ...interface{}) {
	// if !App.WatchEnabled {
	// 	return
	// }

	if watchLock == nil {
		watchLock = new(sync.RWMutex)
	}
	watchLock.Lock()
	defer watchLock.Unlock()

	if watchMap == nil {
		watchMap = make(map[string]*watchValue)
	}

	l := len(nameValuePairs)
	if l%2 != 0 {
		colorLog("[ERRO] cannot add watch variables without even sized.\n")
		return
	}

	for i := 0; i < l; i += 2 {
		k := nameValuePairs[i]
		v := nameValuePairs[i+1]
		rv := reflect.ValueOf(v)

		if rv.Kind() != reflect.Ptr {
			colorLog("[WARN] cannot watch variable( %s ) because [ %s ]\n",
				k, "it's not passed by pointer")
			continue
		}

		watchMap[fmt.Sprint(k)] = &watchValue{
			Value: rv,
			Note:  "",
		}
	}
}

// func pirntWatchValues() {
// 	buf := new(bytes.Buffer)
// 	for k, v := range watchMap {
// 		printWatchValue(buf, k, "", v)
// 	}
// 	fmt.Print(buf.String())
// }

// func printWatchValue(buf *bytes.Buffer, name, prefix string, wv *watchValue) {
// 	switch wv.Type {
// 	case reflect.String:
// 		fmt.Fprintf(buf, "%s:string:%s\"%v\"\n", name, prefix, wv.Value)
// 	case reflect.Ptr:
// 		if wv.Value.IsNil() {
// 			fmt.Fprintf(buf, "%s:pointer:null\n", name)
// 			return
// 		}

// 		printWatchValue(buf, name, "&", &watchValue{
// 			Type:  wv.Value.Elem().Kind(),
// 			Value: wv.Value.Elem(),
// 		})
// 	}
// }

type watchVar struct {
	Kind, Value, Note string
}

func formatWatchVars() map[string]*watchVar {
	watchLock.RLock()
	defer watchLock.RUnlock()

	watchVars := make(map[string]*watchVar)
	for k, v := range watchMap {
		watchVars[k] = &watchVar{
			Kind:  fmt.Sprint(v.Value.Elem().Kind()),
			Value: reflectToStr(v.Value.Elem()),
			Note:  v.Note,
		}
	}
	return watchVars
}

func reflectToStr(rv reflect.Value) string {
	k := rv.Kind()
	switch k {
	case reflect.String:
		return fmt.Sprintf("\"%s\"", rv.String())
	case reflect.Bool:
		return fmt.Sprintf("\"%v\"", rv.Bool())
	default:
		return "Unknown"
	}
}
