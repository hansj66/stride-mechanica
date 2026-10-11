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
	"fmt"
	"mechanica/robot"
	"strings"
	"testing"
	"testing/fstest"
)

// memoryPodStore keeps pods in memory, like the browser version's local storage
type memoryPodStore map[string][]byte

func (m memoryPodStore) SavePod(name string, data []byte) error {
	m[name] = data
	return nil
}

func (m memoryPodStore) LoadPod(name string) ([]byte, error) {
	data, ok := m[name]
	if !ok {
		return nil, fmt.Errorf("no pod called %s", name)
	}
	return data, nil
}

func (m memoryPodStore) PodLocation(name string) string {
	return name + " (in memory)"
}

func TestScriptsFromAFileSystem(t *testing.T) {
	s := newTestShell(t, robot.NewExampleHexapod2())
	s.scripts = fstest.MapFS{
		"square.goik": {Data: []byte("walk 0 50 0 for 2\nwalk 50 0 0 for 2\nhalt\n")},
		"broken.goik": {Data: []byte("walk 0 50 0\njump 3\n")},
		"notes.txt":   {Data: []byte("not a script")},
	}

	if out := s.mustRun("scripts"); !strings.Contains(out, "\tbroken\n\tsquare\n") || strings.Contains(out, "notes") {
		t.Errorf("scripts listed:\n%s", out)
	}
	if out := s.mustRun("run square"); !strings.Contains(out, "Running square (3 commands)") {
		t.Errorf("run square:\n%s", out)
	}
	if _, err := s.run("run broken"); err == nil || !strings.Contains(err.Error(), "broken.goik: line 2") {
		t.Errorf("a broken script gave %v, want the file and line", err)
	}
	for _, name := range []string{"nosuch", "../square"} {
		if _, err := s.run("run " + name); err == nil || !strings.Contains(err.Error(), "there is no script called") {
			t.Errorf("run %s gave %v", name, err)
		}
	}
}

func TestSaveAndLoadPods(t *testing.T) {
	s := newTestShell(t, robot.NewExampleHexapodAX12())
	s.pods = memoryPodStore{}

	s.mustRun("set_coxa_length all 77")
	if out := s.mustRun("save long"); !strings.Contains(out, "Saved long (in memory)") {
		t.Errorf("save:\n%s", out)
	}
	s.mustRun("reset 0")
	s.mustRun("load long")
	if c := s.Pod.BodyDefinition.Segments; c[0].Coxa != 77 || c[5].Coxa != 77 {
		t.Errorf("segment lengths %v after loading, want coxas of 77", c)
	}
	if _, err := s.run("load nothere"); err == nil {
		t.Errorf("loading a pod that was never saved worked")
	}
}
