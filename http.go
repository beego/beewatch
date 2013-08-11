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
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	"code.google.com/p/go.net/websocket"
)

var (
	viewPath string
	data     map[interface{}]interface{}
)

func initHTTP() {
	// Static path.
	sp, err := getStaticPath()
	if err != nil {
		colorLog("[ERRO] BW: Fail to get static path[ %s ]\n", err)
		os.Exit(2)
	}

	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(sp))))
	colorLog("[INIT] BW: File server( %s )\n",
		sp[:strings.LastIndex(sp, "/")])

	// View path.
	viewPath = strings.Replace(sp, "static", "views", 1)

	http.HandleFunc("/", mainPage)
	http.HandleFunc("/gosource", gosource)
	http.Handle("/beewatch", websocket.Handler(connectHandler))

	data = make(map[interface{}]interface{})
	data["AppVer"] = "v" + APP_VER
	data["AppName"] = App.Name

	go sendLoop()
	go listen()
}

func listen() {
	err := http.ListenAndServe(":"+fmt.Sprintf("%d", App.HttpPort), nil)
	if err != nil {
		colorLog("[ERRO] BW: Server crashed[ %s ]\n", err)
		os.Exit(2)
	}
}

func Close() {
	if App.WatchEnabled && !App.CmdMode {
		channelExchangeCommands(LevelCritical, command{Action: "DONE"})
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	b, err := loadFile(viewPath + "/home.html")
	if err != nil {
		colorLog("[ERRO] BW: Fail to load template[ %s ]\n", err)
		os.Exit(2)
	}
	t := template.New("home.html")
	t.Parse(string(b))
	t.Execute(w, data)
}

// serve a (source) file for displaying in the debugger
func gosource(w http.ResponseWriter, r *http.Request) {
	if App.PrintSource {
		fileName := r.FormValue("file")
		// should check for permission?
		w.Header().Set("Cache-control", "no-store, no-cache, must-revalidate")
		http.ServeFile(w, r, fileName)
	} else {
		io.WriteString(w, "'print_source' disenabled.")
	}
}

var (
	currentWebsocket   *websocket.Conn
	connectChannel     = make(chan command)
	toBrowserChannel   = make(chan command)
	fromBrowserChannel = make(chan command)
)

// command is used to transport message to and from the debugger.
type command struct {
	Action     string
	Level      string
	Parameters map[string]string
}

// addParam adds a key,value string pair to the command ; no check on overwrites.
func (c *command) addParam(key, value string) {
	if c.Parameters == nil {
		c.Parameters = map[string]string{}
	}
	c.Parameters[key] = value
}

func connectHandler(ws *websocket.Conn) {
	if currentWebsocket != nil {
		colorLog("[INFO] BW: Connection has already been established, ignore.\n")
		return
	}

	var cmd command
	if err := websocket.JSON.Receive(ws, &cmd); err != nil {
		colorLog("[ERRO] BW: Fail to establish connection[ %s ]\n", err)
	} else {
		currentWebsocket = ws
		connectChannel <- cmd
		App.WatchEnabled = true
		colorLog("[SUCC] BW: Connected to browser, ready to watch.\n")
		receiveLoop()
	}
}

// receiveLoop reads commands from the websocket and puts them onto a channel.
func receiveLoop() {
	for {
		var cmd command
		if err := websocket.JSON.Receive(currentWebsocket, &cmd); err != nil {
			colorLog("[ERRO] BW: connectHandler.JSON.Receive failed[ %s ]\n", err)
			cmd = command{Action: "QUIT"}
		}

		colorLog("[SUCC] BW: Received %v.\n", cmd)
		if "QUIT" == cmd.Action {
			App.WatchEnabled = false
			colorLog("[INFO] BW: Browser requests disconnect.\n")
			currentWebsocket.Close()
			currentWebsocket = nil
			//toBrowserChannel <- cmd
			if cmd.Parameters["PASSIVE"] == "1" {
				fromBrowserChannel <- cmd
			}
			close(toBrowserChannel)
			close(fromBrowserChannel)
			App.WatchEnabled = false
			colorLog("[WARN] BW: Disconnected.\n")
			break
		} else {
			fromBrowserChannel <- cmd
		}
	}

	colorLog("[WARN] BW: Exit receive loop.\n")
	//go sendLoop()
}

// sendLoop takes commands from a channel to send to the browser (debugger).
// If no connection is available then wait for it.
// If the command action is quit then abort the loop.
func sendLoop() {
	if currentWebsocket == nil {
		colorLog("[INFO] BW: No connection, wait for it.\n")
		cmd := <-connectChannel
		if "QUIT" == cmd.Action {
			return
		}
	}

	for {
		next, ok := <-toBrowserChannel
		if !ok {
			colorLog("[WARN] BW: Send channel was closed.\n")
			break
		}

		if "QUIT" == next.Action {
			break
		}

		if currentWebsocket == nil {
			colorLog("[INFO] BW: No connection, wait for it.\n")
			cmd := <-connectChannel
			if "QUIT" == cmd.Action {
				break
			}
		}
		websocket.JSON.Send(currentWebsocket, &next)
		colorLog("[SUCC] BW: Sent %v.\n", next)
	}

	colorLog("[WARN] BW: Exit send loop.\n")
}
