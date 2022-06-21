// Copyright © 2022 Shyam Jeedigunta <shyam123.jvs95@gmail.com>
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

package util

import (
	"encoding/csv"
	"os"
	"sync"

	"k8s.io/klog/v2"
)

type ThreadSafeCsvWriter struct {
	lock      sync.Mutex
	csvWriter *csv.Writer
}

func NewThreadSafeCsvWriter(fileName string) *ThreadSafeCsvWriter {
	csvFile, err := os.Create(fileName)
	if err != nil {
		klog.Errorf("Failed to create file: %v", err)
		os.Exit(1)
	}
	w := csv.NewWriter(csvFile)
	return &ThreadSafeCsvWriter{csvWriter: w}
}

func (w *ThreadSafeCsvWriter) Write(row []string) {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.csvWriter.Write(row)
}

func (w *ThreadSafeCsvWriter) Flush() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.csvWriter.Flush()
}
