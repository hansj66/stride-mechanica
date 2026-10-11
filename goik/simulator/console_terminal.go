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

//go:build !js

package simulator

import "github.com/borud/chatui"

// The status shown when the shell starts
const START_STATUS = "type /quit to exit"

// newConsole returns a chat UI in the terminal the simulator was started from
func newConsole(output chan string, commands chan string) console {
	return chatui.New(chatui.Config{
		OutputCh:     output,
		CommandCh:    commands,
		DynamicColor: false,
		BlockCtrlC:   true,
		HistorySize:  10,
	})
}
