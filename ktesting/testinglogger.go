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
// Experimental
//
// Notice: This package is EXPERIMENTAL and may be changed or removed in a
// later release.
package ktesting

import (
	"bytes"

	"github.com/go-logr/logr"

	"k8s.io/klog/v2/internal/serialize"
	"k8s.io/klog/v2/internal/verbosity"
)

// TL is the relevant subset of testing.TB.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type TL interface {
	Helper()
	Log(args ...interface{})
}

// NewLogger constructs a new logger for the given test interface.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
func NewLogger(t TL, c *Config) logr.Logger {
	return logr.New(&tlogger{
		t:      t,
		prefix: "",
		values: nil,
		config: c,
	})
}

type tlogger struct {
	t      TL
	prefix string
	values []interface{}
	config *Config
}

func (l *tlogger) Init(info logr.RuntimeInfo) {
}

func (l *tlogger) GetCallStackHelper() func() {
	return l.t.Helper
}

func (l *tlogger) Info(level int, msg string, kvList ...interface{}) {
	l.t.Helper()
	buffer := &bytes.Buffer{}
	trimmed := serialize.TrimDuplicates(l.values, kvList)
	serialize.KVListFormat(buffer, trimmed[0]...)
	serialize.KVListFormat(buffer, trimmed[1]...)
	l.log("INFO", msg, buffer)
}

func (l *tlogger) Enabled(level int) bool {
	return l.config.vstate.Enabled(verbosity.Level(level), 1)
}

func (l *tlogger) Error(err error, msg string, kvList ...interface{}) {
	l.t.Helper()
	buffer := &bytes.Buffer{}
	if err != nil {
		serialize.KVListFormat(buffer, "err", err)
	}
	trimmed := serialize.TrimDuplicates(l.values, kvList)
	serialize.KVListFormat(buffer, trimmed[0]...)
	serialize.KVListFormat(buffer, trimmed[1]...)
	l.log("ERROR", msg, buffer)
}

func (l *tlogger) log(what, msg string, buffer *bytes.Buffer) {
	l.t.Helper()
	args := []interface{}{what}
	if l.prefix != "" {
		args = append(args, l.prefix+":")
	}
	args = append(args, msg)
	if buffer.Len() > 0 {
		// Skip leading space inserted by serialize.KVListFormat.
		args = append(args, string(buffer.Bytes()[1:]))
	}
	l.t.Log(args...)
}

// WithName returns a new logr.Logger with the specified name appended.  klogr
// uses '/' characters to separate name elements.  Callers should not pass '/'
// in the provided name string, but this library does not actually enforce that.
func (l *tlogger) WithName(name string) logr.LogSink {
	new := *l
	if len(l.prefix) > 0 {
		new.prefix = l.prefix + "/"
	}
	new.prefix += name
	return &new
}

func (l *tlogger) WithValues(kvList ...interface{}) logr.LogSink {
	new := *l
	new.values = serialize.WithValues(l.values, kvList)
	return &new
}

var _ logr.LogSink = &tlogger{}
var _ logr.CallStackHelperLogSink = &tlogger{}
