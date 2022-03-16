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

package klog_test

import (
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

func ExampleSetLogger() {
	defer klog.ClearLogger()

	// Logger is only used as backend, Background() returns klogr.
	klog.SetLogger(logr.Discard())
	fmt.Printf("logger after SetLogger: %T\n", klog.Background().GetSink())

	// Logger is only used as backend, Background() returns klogr.
	klog.SetLoggerWithOptions(logr.Discard(), klog.ContextualLogger(false))
	fmt.Printf("logger after SetLoggerWithOptions with ContextualLogger(false): %T\n", klog.Background().GetSink())

	// Logger is used as backend and directly.
	klog.SetLoggerWithOptions(logr.Discard(), klog.ContextualLogger(true))
	fmt.Printf("logger after SetLoggerWithOptions with ContextualLogger(true): %T\n", klog.Background().GetSink())

	// Output:
	// logger after SetLogger: *klog.klogger
	// logger after SetLoggerWithOptions with ContextualLogger(false): *klog.klogger
	// logger after SetLoggerWithOptions with ContextualLogger(true): logr.discardLogSink
}

func ExampleFlushLogger() {
	defer klog.ClearLogger()

	// This simple logger doesn't need flushing, but others might.
	klog.SetLoggerWithOptions(logr.Discard(), klog.FlushLogger(func() {
		fmt.Print("flushing...")
	}))
	klog.Flush()

	// Output:
	// flushing...
}
