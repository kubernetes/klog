/*
Copyright 2019 The Kubernetes Authors.
Copyright 2020 Intel Coporation.

SPDX-License-Identifier: Apache-2.0
*/

package testinglogger

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func TestInfo(t *testing.T) {
	tests := map[string]struct {
		text           string
		withValues     []interface{}
		keysAndValues  []interface{}
		names          []string
		err            error
		expectedOutput string
	}{
		"should log with values passed to keysAndValues": {
			text:          "test",
			keysAndValues: []interface{}{"akey", "avalue"},
			expectedOutput: `INFO test akey="avalue"
`,
		},
		"should support single name": {
			names:         []string{"hello"},
			text:          "test",
			keysAndValues: []interface{}{"akey", "avalue"},
			expectedOutput: `INFO hello: test akey="avalue"
`,
		},
		"should support multiple names": {
			names:         []string{"hello", "world"},
			text:          "test",
			keysAndValues: []interface{}{"akey", "avalue"},
			expectedOutput: `INFO hello/world: test akey="avalue"
`,
		},
		"should not print duplicate keys with the same value": {
			text:          "test",
			keysAndValues: []interface{}{"akey", "avalue", "akey", "avalue"},
			expectedOutput: `INFO test akey="avalue"
`,
		},
		"should only print the last duplicate key when the values are passed to Info": {
			text:          "test",
			keysAndValues: []interface{}{"akey", "avalue", "akey", "avalue2"},
			expectedOutput: `INFO test akey="avalue2"
`,
		},
		"should only print the duplicate key that is passed to Info if one was passed to the logger": {
			withValues:    []interface{}{"akey", "avalue"},
			text:          "test",
			keysAndValues: []interface{}{"akey", "avalue"},
			expectedOutput: `INFO test akey="avalue"
`,
		},
		"should only print the key passed to Info when one is already set on the logger": {
			withValues:    []interface{}{"akey", "avalue"},
			text:          "test",
			keysAndValues: []interface{}{"akey", "avalue2"},
			expectedOutput: `INFO test akey="avalue2"
`,
		},
		"should correctly handle odd-numbers of KVs": {
			text:          "test",
			keysAndValues: []interface{}{"akey", "avalue", "akey2"},
			expectedOutput: `INFO test akey="avalue" akey2=<nil>
`,
		},
		"should correctly html characters": {
			text:          "test",
			keysAndValues: []interface{}{"akey", "<&>"},
			expectedOutput: `INFO test akey="<&>"
`,
		},
		"should correctly handle odd-numbers of KVs in both log values and Info args": {
			withValues:    []interface{}{"basekey1", "basevar1", "basekey2"},
			text:          "test",
			keysAndValues: []interface{}{"akey", "avalue", "akey2"},
			expectedOutput: `INFO test basekey1="basevar1" basekey2=<nil> akey="avalue" akey2=<nil>
`,
		},
		"should correctly print regular error types": {
			text:          "test",
			keysAndValues: []interface{}{"err", errors.New("whoops")},
			expectedOutput: `INFO test err="whoops"
`,
		},
		"should correctly print regular error types when using logr.Error": {
			text: "test",
			err:  errors.New("whoops"),
			expectedOutput: `ERROR test err="whoops"
`,
		},
	}
	for n, test := range tests {
		t.Run(n, func(t *testing.T) {
			var buffer logToBuf
			klogr := New(&buffer)
			for _, name := range test.names {
				klogr = klogr.WithName(name)
			}
			klogr = klogr.WithValues(test.withValues...)

			if test.err != nil {
				klogr.Error(test.err, test.text, test.keysAndValues...)
			} else {
				klogr.Info(test.text, test.keysAndValues...)
			}

			actual := buffer.String()
			if actual != test.expectedOutput {
				t.Errorf("expected %q did not match actual %q", test.expectedOutput, actual)
			}
		})
	}
}

type logToBuf struct {
	bytes.Buffer
}

func (l *logToBuf) Helper() {
}

func (l *logToBuf) Log(args ...interface{}) {
	for i, arg := range args {
		if i > 0 {
			l.Write([]byte(" "))
		}
		l.Write([]byte(fmt.Sprintf("%s", arg)))
	}
	l.Write([]byte("\n"))
}
