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
)

// CAD export (see docs/cad-export.md and cad/goik_cad.py). The assembly description is written to cad/out and
// built with CadQuery on the desktop (cadOutput_desktop.go), and downloaded in the browser version (cadOutput_js.go)

func (s *Shell) executeExportCadCmd(args []string) error {
	s.outputCh <- fmt.Sprintf("%+v", args)

	if len(args) != 2 {
		return fmt.Errorf("syntax error ('export_cad <name>'): %+v", args)
	}
	name := args[1]
	if strings.ContainsAny(name, `/\:`) || strings.HasPrefix(name, ".") {
		return fmt.Errorf("the name can not contain a path")
	}

	assembly, err := robot.NewCadAssembly(s.Pod.BodyDefinition, name)
	if err != nil {
		return err
	}
	return s.outputCad(assembly, name)
}
