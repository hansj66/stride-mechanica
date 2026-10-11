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

package main

import (
	"embed"
	"io/fs"
	"log"
	"mechanica/simulator"
)

// The browser version (see cmd/web and web/index.html) can't read the scripts folder, so the scripts are built in
//
//go:embed scripts/*.goik
var builtInScripts embed.FS

func main() {
	scripts, err := fs.Sub(builtInScripts, "scripts")
	if err != nil {
		log.Fatal(err)
	}
	simulator.Run(simulator.Options{Scripts: scripts})
}
