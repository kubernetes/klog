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
package klog

import (
	"fmt"
	"runtime"
)

// numTotalFrames specifies the maximum number of frames that are supported by
// backtrace APIs.
const numTotalFrames = 100

type traceT struct {
	buf       []byte
	all       bool
	skip      int
	numframes int
}

func (t traceT) String() string {
	return string(t.buf)
}

// BacktraceOption provides functional parameters for Backtrace.
type BacktraceOption func(*traceT)

// BacktraceAll returns backtraces for all goroutines.
// Note: This is incompatible with other options like skip/size.
// Those will be ignored when tracing all goroutines.
func BacktraceAll(all bool) BacktraceOption {
	return func(t *traceT) {
		t.all = all
	}
}

// BacktraceSkip is like Backtrace except that it skips some stack levels.
// BacktraceSkip(0) is equivalent to Backtrace(). Has no effect when
// combined with BacktraceAll(true).
func BacktraceSkip(depth int) BacktraceOption {
	return func(t *traceT) {
		t.skip = depth
	}
}

// BacktraceSize will limit how far the unwinding goes, i.e. specify
// how many stack frames will be printed. Has no effect when
// combined with BacktraceAll(true).
func BacktraceSize(numFrames int) BacktraceOption {
	return func(t *traceT) {
		if numFrames > 0 {
			t.numframes = numFrames
		}
	}
}

// Backtrace returns an object that as default represents the stack backtrace of the calling
// goroutine. That object can be used as value in a structured logging call.
// It supports printing as string or as structured output via logr.MarshalLog.
// The behavior can be modified via options.
func Backtrace(opts ...BacktraceOption) interface{} {

	trace := traceT{skip: 0, numframes: 0}

	for _, opt := range opts {
		opt(&trace)
	}

	// 'All' supersedes skip/size etc
	if trace.all {
		trace.buf = stacks(true)
		return trace
	}

	pc := make([]uintptr, numTotalFrames)
	// skip runtime.Callers and the klog.Backtrace API
	n := runtime.Callers(trace.skip+2, pc)

	if n == 0 {
		// No PCs available. This can happen if the first argument to
		// runtime.Callers is large.
		//
		// Return now to avoid processing the zero Frame that would
		// otherwise be returned by frames.Next below.
		return nil
	}

	if n > numTotalFrames {
		fmt.Printf("error: runtime.Callers returned too many pcs (>%v)\n", numTotalFrames)
		return nil
	}

	// pass only valid pcs to runtime.CallersFrames (remove goexit..)
	pc = pc[:n-1]

	// Account for "size" parameter
	if trace.numframes > 0 && trace.numframes < n {
		pc = pc[:trace.numframes]
	}

	frames := runtime.CallersFrames(pc)

	var s string
	for {
		frame, more := frames.Next()
		s = s + fmt.Sprintf("%s():\n\t%s:%v\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	if s != "" {
		trace.buf = []byte(s)
	}

	return trace
}
