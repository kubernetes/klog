/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package klogr_test

import (
	"bytes"
	"flag"
	"strings"
	"testing"

	"k8s.io/klog/v2"
)

// TestSetKLogLogger checks that it is possible to use a modified KLog Logger
// as contextual logger.
func TestSetKLogLogger(t *testing.T) {
	defer klog.CaptureState().Restore()
	buf := new(bytes.Buffer)
	klog.SetOutput(buf)
	var fs flag.FlagSet
	klog.InitFlags(&fs)
	fs.Set("logtostderr", "false")
	logger := klog.Background()
	logger = logger.WithValues("hello", "world")
	klog.SetLoggerWithOptions(logger, klog.ContextualLogger(true))
	logger2 := klog.Background()

	if logger != logger2 {
		t.Fatalf("Expected to get the modified logger %+v from Background, got: %+v", logger, logger2)
	}

	logger2.Info("ping")

	str := buf.String()
	expected := `"ping" hello="world"`
	if !strings.Contains(str, expected) {
		t.Fatalf("Expected %q to contain %q", str, expected)
	}
}
