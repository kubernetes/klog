/*
Copyright 2019 The Kubernetes Authors.
Copyright 2020 Intel Coporation.

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

// Package testinglogger contains an implementation of the logr interface
// which is logging through a function like testing.TB.Log function.
// Therefore it can be used in standard Go tests and Gingko test suites
// to ensure that output is associated with the currently running test.
//
// Serialization of the structured log parameters is done in the same way
// as for klog.InfoS.
//
package testinglogger

import (
	"bytes"

	"github.com/go-logr/logr"

	"k8s.io/klog/v2/internal/serialize"
)

// Logger is the relevant subset of testing.TB.
type Logger interface {
	Helper()
	Log(args ...interface{})
}

// New constructs a new logger for the given test interface.
func New(l Logger) logr.Logger {
	return logr.New(&tlogger{
		l:      l,
		prefix: "",
		values: nil,
	})
}

type tlogger struct {
	l      Logger
	prefix string
	values []interface{}
}

func (l *tlogger) clone() *tlogger {
	return &tlogger{
		l:      l.l,
		prefix: l.prefix,
		values: copySlice(l.values),
	}
}

func copySlice(in []interface{}) []interface{} {
	out := make([]interface{}, len(in))
	copy(out, in)
	return out
}

func (l *tlogger) Init(info logr.RuntimeInfo) {
}

func (l *tlogger) GetCallStackHelper() func() {
	return l.l.Helper
}

func (l *tlogger) Info(level int, msg string, kvList ...interface{}) {
	l.l.Helper()
	buffer := &bytes.Buffer{}
	trimmed := serialize.TrimDuplicates(l.values, kvList)
	serialize.KVListFormat(buffer, trimmed[0]...)
	serialize.KVListFormat(buffer, trimmed[1]...)
	l.log("INFO", msg, buffer)
}

func (l *tlogger) Enabled(level int) bool {
	return true
}

func (l *tlogger) Error(err error, msg string, kvList ...interface{}) {
	l.l.Helper()
	buffer := &bytes.Buffer{}
	serialize.KVListFormat(buffer, "err", err)
	trimmed := serialize.TrimDuplicates(l.values, kvList)
	serialize.KVListFormat(buffer, trimmed[0]...)
	serialize.KVListFormat(buffer, trimmed[1]...)
	l.log("ERROR", msg, buffer)
}

func (l *tlogger) log(what, msg string, buffer *bytes.Buffer) {
	l.l.Helper()
	args := []interface{}{what}
	if l.prefix != "" {
		args = append(args, l.prefix+":")
	}
	args = append(args, msg)
	if buffer.Len() > 0 {
		// Skip leading space inserted by serialize.KVListFormat.
		args = append(args, string(buffer.Bytes()[1:]))
	}
	l.l.Log(args...)
}

// WithName returns a new logr.Logger with the specified name appended.  klogr
// uses '/' characters to separate name elements.  Callers should not pass '/'
// in the provided name string, but this library does not actually enforce that.
func (l *tlogger) WithName(name string) logr.LogSink {
	new := l.clone()
	if len(l.prefix) > 0 {
		new.prefix = l.prefix + "/"
	}
	new.prefix += name
	return new
}

func (l *tlogger) WithValues(kvList ...interface{}) logr.LogSink {
	new := l.clone()
	new.values = append(new.values, kvList...)
	return new
}

var _ logr.LogSink = &tlogger{}
var _ logr.CallStackHelperLogSink = &tlogger{}
