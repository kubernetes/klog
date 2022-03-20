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

// Package textlogger contains an implementation of the logr interface
// which is producing the exact same output as klog.
//
// Experimental
//
// Notice: This package is EXPERIMENTAL and may be changed or removed in a
// later release.
package textlogger

import (
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"

	"k8s.io/klog/v2/internal/buffer"
	"k8s.io/klog/v2/internal/serialize"
	"k8s.io/klog/v2/internal/severity"
	"k8s.io/klog/v2/internal/verbosity"
)

var (
	// TimeNow is used to retrieve the current time. May be changed for testing.
	//
	// Experimental
	//
	// Notice: This variable is EXPERIMENTAL and may be changed or removed in a
	// later release.
	TimeNow = time.Now
)

// NewLogger constructs a new logger.
//
// Experimental
//
// Notice: This function is EXPERIMENTAL and may be changed or removed in a
// later release. The behavior of the returned Logger may change.
func NewLogger(c *Config) logr.Logger {
	return logr.New(&tlogger{
		prefix:      "",
		values:      nil,
		config:      c,
		bufferCache: &buffer.Buffers{},
	})
}

type tlogger struct {
	callDepth   int
	prefix      string
	values      []interface{}
	config      *Config
	bufferCache *buffer.Buffers
}

func copySlice(in []interface{}) []interface{} {
	out := make([]interface{}, len(in))
	copy(out, in)
	return out
}

func (l *tlogger) Init(info logr.RuntimeInfo) {
	l.callDepth = info.CallDepth
}

func (l *tlogger) WithCallDepth(depth int) logr.LogSink {
	newLogger := *l
	newLogger.callDepth += depth
	return &newLogger
}

func (l *tlogger) Enabled(level int) bool {
	return l.config.Enabled(verbosity.Level(level), 1)
}

func (l *tlogger) Info(level int, msg string, kvList ...interface{}) {
	l.print(nil, severity.InfoLog, msg, kvList)
}

func (l *tlogger) Error(err error, msg string, kvList ...interface{}) {
	l.print(err, severity.ErrorLog, msg, kvList)
}

func (l *tlogger) print(err error, s severity.Severity, msg string, kvList []interface{}) {
	// Only create a new buffer if we don't have one cached.
	b := l.bufferCache.GetBuffer()

	// Determine caller.
	// +1 for this frame, +1 for Info/Error.
	_, file, line, ok := runtime.Caller(l.callDepth + 2)
	if !ok {
		file = "???"
		line = 1
	} else {
		if slash := strings.LastIndex(file, "/"); slash >= 0 {
			path := file
			file = path[slash+1:]
		}
	}

	// Format header.
	now := TimeNow()
	b.FormatHeader(s, file, line, now)

	// Inject WithName names into message.
	if l.prefix != "" {
		msg = l.prefix + ": " + msg
	}

	// The message is always quoted, even if it contains line breaks.
	// If developers want multi-line output, they should use a small, fixed
	// message and put the multi-line output into a value.
	b.WriteString(strconv.Quote(msg))
	if err != nil {
		serialize.KVListFormat(&b.Buffer, "err", err)
	}
	trimmed := serialize.TrimDuplicates(l.values, kvList)
	serialize.KVListFormat(&b.Buffer, trimmed[0]...)
	serialize.KVListFormat(&b.Buffer, trimmed[1]...)
	if b.Len() == 0 || b.Bytes()[b.Len()-1] != '\n' {
		b.WriteByte('\n')
	}
	l.config.co.output.Write(b.Bytes())
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
var _ logr.CallDepthLogSink = &tlogger{}
