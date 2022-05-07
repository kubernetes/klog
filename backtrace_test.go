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
	"regexp"
	"runtime"
	"strings"
	"testing"

	"k8s.io/klog/v2"
)

// The tests are sensitive to line changes in the following code. Here are some
// lines that can be added or removed to compensate for import statement
// changes.
//
//
//
//
//
//
//
//
// This must be line 40 (first line is 1).

func outer(callback func() interface{}) interface{} {
	return inner(callback)
}

func inner(callback func() interface{}) interface{} {
	return callback()
}

//
//
//
//
//
//
//
//
//
//
// This must be line 60. Any line number higher than that gets ignored by normalizeBacktrace.

func TestBacktrace(t *testing.T) {
	for name, tt := range map[string]struct {
		callback func() interface{}
		expected string
	}{
		"simple": {
			callback: func() interface{} { return klog.Backtrace() },
			expected: `k8s.io/klog/v2_test.TestBacktrace.funcX():
	/zzz/backtrace_test.go:xxx
k8s.io/klog/v2_test.inner():
	/zzz/backtrace_test.go:47
k8s.io/klog/v2_test.outer():
	/zzz/backtrace_test.go:43
k8s.io/klog/v2_test.TestBacktrace.funcX():
	/zzz/backtrace_test.go:xxx
`,
		},

		"skip0": {
			callback: func() interface{} { return klog.Backtrace(klog.BacktraceSkip(0)) },
			expected: `k8s.io/klog/v2_test.TestBacktrace.funcX():
	/zzz/backtrace_test.go:xxx
k8s.io/klog/v2_test.inner():
	/zzz/backtrace_test.go:47
k8s.io/klog/v2_test.outer():
	/zzz/backtrace_test.go:43
k8s.io/klog/v2_test.TestBacktrace.funcX():
	/zzz/backtrace_test.go:xxx
`,
		},

		"skip1": {
			callback: func() interface{} { return klog.Backtrace(klog.BacktraceSkip(1)) },
			expected: `k8s.io/klog/v2_test.inner():
	/zzz/backtrace_test.go:47
k8s.io/klog/v2_test.outer():
	/zzz/backtrace_test.go:43
k8s.io/klog/v2_test.TestBacktrace.funcX():
	/zzz/backtrace_test.go:xxx
`,
		},

		"skip2": {
			callback: func() interface{} { return klog.Backtrace(klog.BacktraceSkip(2)) },
			expected: `k8s.io/klog/v2_test.outer():
	/zzz/backtrace_test.go:43
k8s.io/klog/v2_test.TestBacktrace.funcX():
	/zzz/backtrace_test.go:xxx
`,
		},

		"trace1": {
			callback: func() interface{} { return klog.Backtrace(klog.BacktraceSize(1)) },
			expected: `k8s.io/klog/v2_test.TestBacktrace.funcX():
	/zzz/backtrace_test.go:xxx
`,
		},

		"trace2": {
			callback: func() interface{} { return klog.Backtrace(klog.BacktraceSize(2)) },
			expected: `k8s.io/klog/v2_test.TestBacktrace.funcX():
	/zzz/backtrace_test.go:xxx
k8s.io/klog/v2_test.inner():
	/zzz/backtrace_test.go:47
`,
		},

		"skip1-trace1": {
			callback: func() interface{} { return klog.Backtrace(klog.BacktraceSize(1), klog.BacktraceSkip(1)) },
			expected: `k8s.io/klog/v2_test.inner():
	/zzz/backtrace_test.go:47
`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			backtrace := outer(tt.callback)
			actual := normalizeBacktrace(t, backtrace)
			if actual != tt.expected {
				t.Errorf("Wrong backtrace. Expected:\n%s\nActual:\n%s\n", tt.expected, actual)
			}
		})
	}
}

func normalizeBacktrace(t *testing.T, backtrace interface{}) string {
	str := backtrace.(fmt.Stringer).String()

	// This matches the stack entry for the testing package:
	// testing.tRunner():
	//     /nvme/gopath/go-1.18.1/src/testing/testing.go:1439
	//
	// It and all following entries vary and get removed.
	end := regexp.MustCompile(`(?m)^testing\.`).FindStringIndex(str)
	if end != nil {
		str = str[:end[0]]
	}

	// Remove varying line numbers.
	str = regexp.MustCompile(`(?m):([67890][0-9]|[1-9][0-9][0-9])$`).ReplaceAllString(str, ":xxx")

	// Remove varying anonymous function numbers.
	str = regexp.MustCompile(`\.func[[:digit:]]+`).ReplaceAllString(str, ".funcX")

	// Remove varying path names
	str = regexp.MustCompile(`([[:blank:]]+)[[:^blank:]]*(backtrace_test.go:.*)`).ReplaceAllString(str, "$1/zzz/$2")

	return str
}

func TestBacktraceAll(t *testing.T) {
	_, callerFile, _, _ := runtime.Caller(0)
	actual := klog.Backtrace(klog.BacktraceAll(true)).(fmt.Stringer).String()
	t.Logf("Backtrace:\n%s\n", actual)

	if strings.Contains(actual, callerFile) == false {
		t.Errorf("Expected to see %q in trace:\n%s", callerFile, actual)
	}

	// Pattern: goroutine 7 [running]:
	p := regexp.MustCompile("(?m)(^goroutine [0-9]*.\\[.*\\]*:)")
	if len(p.FindAllString(actual, -1)) < 2 {
		t.Errorf("Expected more than 1 goroutine stack to be printed, got:\n%s", actual)
	}
}
