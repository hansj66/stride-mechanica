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
	"fmt"
	"mechanica/robot"
)

// outputCad downloads the assembly description: CadQuery doesn't run in the browser, so the STEP assembly is built
// from it on a desktop
func (s *Shell) outputCad(assembly *robot.CadAssembly, name string) error {
	data, err := assembly.JSON()
	if err != nil {
		return err
	}
	page := pageConsole()
	if page.IsUndefined() {
		return fmt.Errorf("the page can't download files")
	}
	page.Call("download", name+".json", string(data))
	s.outputCh <- fmt.Sprintf("Downloaded %s.json (%s, rest pose)", name, assembly.String())
	s.outputCh <- "The browser can't build CAD models. Build the STEP assembly and STL from it with the desktop version:"
	s.outputCh <- fmt.Sprintf("\tpython cad/goik_cad.py %s.json (see docs/cad-export.md)", name)
	return nil
}
