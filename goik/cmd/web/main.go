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

// Builds the browser version of Stride Mechanica into the web folder, and serves it on http://localhost:8080.
// Run it from the goik/goik folder:
//
//	go run ./cmd/web           build, then serve the web folder
//	go run ./cmd/web -build    only build (the web folder can then be published as it is, see docs/browser.md)
//	go run ./cmd/web -addr :9000
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const WEB_FOLDER = "web"

func main() {
	buildOnly := flag.Bool("build", false, "only build, don't serve")
	addr := flag.String("addr", "localhost:8080", "address to serve on")
	flag.Parse()

	if err := build(); err != nil {
		log.Fatal(err)
	}
	if *buildOnly {
		return
	}

	url := *addr
	if strings.HasPrefix(url, ":") {
		url = "localhost" + url
	}
	fmt.Printf("Serving %s on http://%s (Ctrl+C to stop)\n", WEB_FOLDER, url)
	files := http.FileServer(http.Dir(WEB_FOLDER))
	log.Fatal(http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always the latest build
		w.Header().Set("Cache-Control", "no-cache")
		files.ServeHTTP(w, r)
	})))
}

// build compiles the simulator to WebAssembly, and copies the Go runtime's JavaScript support next to it. The two
// must come from the same Go version
func build() error {
	fmt.Println("Building", filepath.Join(WEB_FOLDER, "mechanica.wasm"), "...")
	cmd := exec.Command("go", "build", "-trimpath", "-ldflags=-s -w", "-o", filepath.Join(WEB_FOLDER, "mechanica.wasm"), ".")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("building the browser version failed: %w", err)
	}

	out, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		return err
	}
	root := strings.TrimSpace(string(out))
	// lib/wasm since Go 1.24, misc/wasm before
	for _, folder := range []string{"lib", "misc"} {
		support, err := os.ReadFile(filepath.Join(root, folder, "wasm", "wasm_exec.js"))
		if err == nil {
			return os.WriteFile(filepath.Join(WEB_FOLDER, "wasm_exec.js"), support, 0644)
		}
	}
	return fmt.Errorf("wasm_exec.js was not found in %s", root)
}
