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

import (
	"mechanica/robot"
	"testing"
)

func TestCommandNamesInAnyCase(t *testing.T) {
	s := newTestShell(t, robot.NewExampleHexapod2())

	// Phone keyboards capitalize the first letter
	for _, command := range []string{"Walk 0 40 0", "HALT", "Pitch 5", "Speed 10"} {
		if _, err := s.run(command); err != nil {
			t.Errorf("%s: %v", command, err)
		}
	}
	if _, err := s.run("Jump 3"); err == nil {
		t.Errorf("an unknown command was accepted")
	}
}
