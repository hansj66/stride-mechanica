// Copyright 2025 Hans Jørgen Grimstad
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build js

package simulator

import (
	"strings"
	"sync"
	"syscall/js"
)

/*
	The browser version's console

	The console is the page around the views (web/index.html). The views run in a frame (web/game.html), which
	makes the page's mechanicaConsole object available to this code:

	- mechanicaConsole.print(line) shows a line of output, and mechanicaConsole.status(text) a line of status
	- mechanicaConsole.ready() tells the page that commands can be sent: the page then calls mechanicaCommand(text)
	  (defined here, in the frame) for each command typed
	- mechanicaConsole.download(filename, text) saves a file in the browser's downloads (see cadOutput_js.go)
*/

// The status shown when the shell starts
const START_STATUS = "type help for a list of commands"

type webConsole struct {
	output   <-chan string
	commands chan<- string
	stop     chan struct{}
	once     sync.Once
}

func newConsole(output chan string, commands chan string) console {
	return &webConsole{output: output, commands: commands, stop: make(chan struct{})}
}

// pageConsole returns the page's mechanicaConsole object (undefined if the frame was opened on its own)
func pageConsole() js.Value {
	return js.Global().Get("mechanicaConsole")
}

func (c *webConsole) print(line string) {
	if page := pageConsole(); !page.IsUndefined() {
		page.Call("print", line)
	}
}

func (c *webConsole) Run() error {
	command := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		text := args[0].String()
		// This runs on the browser's event loop, which must not block: hand the command over on a goroutine
		go func() {
			if strings.EqualFold(strings.TrimSpace(text), "/quit") {
				c.print("Close the browser tab to quit")
				return
			}
			c.commands <- text
		}()
		return nil
	})
	defer command.Release()
	js.Global().Set("mechanicaCommand", command)
	// Show what is queued already (the banner) before the page can send commands
	for queued := true; queued; {
		select {
		case line := <-c.output:
			c.print(line)
		default:
			queued = false
		}
	}
	if page := pageConsole(); !page.IsUndefined() {
		page.Call("ready")
	}

	for {
		select {
		case line := <-c.output:
			c.print(line)
		case <-c.stop:
			return nil
		}
	}
}

func (c *webConsole) Stop() {
	c.once.Do(func() { close(c.stop) })
}

func (c *webConsole) SetStatus(status string) {
	if page := pageConsole(); !page.IsUndefined() {
		page.Call("status", status)
	}
}
