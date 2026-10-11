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

import (
	"os"
	"path/filepath"
)

const POD_FOLDER = "pods"

// folderPodStore keeps pods as files in a folder
type folderPodStore struct {
	folder string
}

func newPodStore() podStore {
	return folderPodStore{folder: POD_FOLDER}
}

func (f folderPodStore) SavePod(name string, data []byte) error {
	if err := os.MkdirAll(f.folder, 0755); err != nil {
		return err
	}
	return os.WriteFile(f.PodLocation(name), data, 0644)
}

func (f folderPodStore) LoadPod(name string) ([]byte, error) {
	return os.ReadFile(f.PodLocation(name))
}

func (f folderPodStore) PodLocation(name string) string {
	return filepath.Join(f.folder, name)
}
