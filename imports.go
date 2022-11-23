/*
Copyright 2021 The Kubernetes Authors.

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

package klog

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
)

// The reason for providing these aliases is to allow code to work with logr
// without directly importing it.

// Logger in this package is exactly the same as [logr.Logger].
type Logger = logr.Logger

// LogSink in this package is exactly the same as [logr.LogSink].
type LogSink = logr.LogSink

// Runtimeinfo in this package is exactly the same as [logr.RuntimeInfo].
type RuntimeInfo = logr.RuntimeInfo

// Marshaler in this package is exactly the same as [logr.Marshaler].
type Marshaler = logr.Marshaler

// PseudoStruct in this package is exactly the same as [funcr.PseudoStruct].
// Use it as a value to render some key/value pairs like a struct.
type PseudoStruct = funcr.PseudoStruct

var (
	// New is an alias for logr.New.
	New = logr.New
)
