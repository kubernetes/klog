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
	"context"

	"github.com/go-logr/logr"
)

// contextLogSink inherits most of the functionality from some other
// LogSink. It stores a context and information which key/value pairs from that
// context are meant to be logged. Then if some log entry really gets emitted,
// it extracts those values and adds them.
//
// This causes little overhead when nothing gets logged (logger not used or
// entry not enabled).
type contextLogSink struct {
	logr.LogSink
	contextValues
}

// contextLogSinkDepth is a variant of contextLogSink for LogSinks which
// implement also CallDepthLogSink.
type contextLogSinkDepth struct {
	logr.LogSink
	logr.CallDepthLogSink
	contextValues
	callDepth int
}

// contextLogSinkDepth is a variant of contextLogSink for LogSinks which
// implement also CallStackHelperLogSink.
type contextLogSinkHelper struct {
	logr.LogSink
	logr.CallStackHelperLogSink
	contextValues
}

func (cls *contextLogSink) Info(level int, msg string, keysAndValues ...interface{}) {
	keysAndValues = cls.contextValues.append(keysAndValues)
	cls.LogSink.Info(level, msg, keysAndValues...)
}

func (cls *contextLogSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	clone := *cls
	clone.LogSink = clone.LogSink.WithValues(keysAndValues...)
	return &clone
}

func (cls *contextLogSinkDepth) Info(level int, msg string, keysAndValues ...interface{}) {
	keysAndValues = cls.contextValues.append(keysAndValues)
	logSink := cls.CallDepthLogSink.WithCallDepth(1 + cls.callDepth)
	logSink.Info(level, msg, keysAndValues...)
}

func (cls *contextLogSinkDepth) WithCallDepth(depth int) logr.LogSink {
	clone := *cls
	clone.callDepth += depth
	return &clone
}

func (cls *contextLogSinkDepth) WithValues(keysAndValues ...interface{}) logr.LogSink {
	clone := *cls
	clone.LogSink = clone.LogSink.WithValues(keysAndValues...)
	// This expects the same capabilities for the new LogSink as before.
	clone.CallDepthLogSink = clone.LogSink.(logr.CallDepthLogSink)
	return &clone
}

func (cls *contextLogSinkHelper) Info(level int, msg string, keysAndValues ...interface{}) {
	keysAndValues = cls.contextValues.append(keysAndValues)
	helper := cls.CallStackHelperLogSink.GetCallStackHelper()
	helper()
	cls.LogSink.Info(level, msg, keysAndValues...)
}

func (cls *contextLogSinkHelper) WithValues(keysAndValues ...interface{}) logr.LogSink {
	clone := *cls
	clone.LogSink = clone.LogSink.WithValues(keysAndValues...)
	// This expects the same capabilities for the new LogSink as before.
	clone.CallStackHelperLogSink = clone.LogSink.(logr.CallStackHelperLogSink)
	return &clone
}

type contextValues struct {
	ctx  context.Context
	keys []ContextKey
}

func (cv contextValues) append(keysAndValues []interface{}) []interface{} {
	for _, key := range cv.keys {
		if value := cv.ctx.Value(key.Key); value != nil {
			keysAndValues = append(keysAndValues, key.Name, value)
		}
	}
	return keysAndValues
}

func (cv contextValues) extract() []interface{} {
	keysAndValues := make([]interface{}, 0, 2*len(cv.keys))
	for _, key := range cv.keys {
		if value := cv.ctx.Value(key.Key); value != nil {
			keysAndValues = append(keysAndValues, key.Name, value)
		}
	}
	return keysAndValues
}
