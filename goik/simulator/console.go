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

package simulator

// console is the command shell's user interface: a terminal on the desktop (console_terminal.go), or the page
// around the views in the browser version (console_js.go). It shows the lines sent on the shell's output channel,
// and sends the commands typed on its command channel. newConsole creates the platform's console.
type console interface {
	// Run shows the console until it is stopped
	Run() error
	Stop()
	// SetStatus shows a line of status (the last command, for example)
	SetStatus(status string)
}
