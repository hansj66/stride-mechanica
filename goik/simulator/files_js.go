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
	"syscall/js"
)

// Pods are kept in the browser's local storage under this prefix and their name
const POD_STORAGE_PREFIX = "mechanica/pods/"

// localPodStore keeps pods in the browser's local storage: they stay in this browser, on this computer
type localPodStore struct{}

func newPodStore() podStore {
	return localPodStore{}
}

func localStorage() (js.Value, error) {
	storage := js.Global().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return js.Value{}, fmt.Errorf("this browser has no local storage to keep pods in")
	}
	return storage, nil
}

func (localPodStore) SavePod(name string, data []byte) (err error) {
	storage, err := localStorage()
	if err != nil {
		return err
	}
	// setItem throws when the storage is full or blocked
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("the browser refused to store the pod: %v", r)
		}
	}()
	storage.Call("setItem", POD_STORAGE_PREFIX+name, string(data))
	return nil
}

func (localPodStore) LoadPod(name string) ([]byte, error) {
	storage, err := localStorage()
	if err != nil {
		return nil, err
	}
	item := storage.Call("getItem", POD_STORAGE_PREFIX+name)
	if item.IsNull() {
		return nil, fmt.Errorf("no pod called %s is saved in this browser", name)
	}
	return []byte(item.String()), nil
}

func (localPodStore) PodLocation(name string) string {
	return fmt.Sprintf("%s (in this browser's local storage)", name)
}
