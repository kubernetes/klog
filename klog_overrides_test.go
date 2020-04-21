// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package klog

import (
	"testing"
)

type klogWriter struct{}

func (kw klogWriter) Write(p []byte) (n int, err error) {
	// OutputCallDepth is the depth in the stack where the call to klog methods
	// originate from. DefaultPrefixLength is the length of the prefix before the
	// actual thing we need to log
	if p[0] == 'I' {
		InfoDepth(OutputCallDepth, string(p[DefaultPrefixLength:]))
	} else if p[0] == 'W' {
		WarningDepth(OutputCallDepth, string(p[DefaultPrefixLength:]))
	} else if p[0] == 'E' {
		ErrorDepth(OutputCallDepth, string(p[DefaultPrefixLength:]))
	} else if p[0] == 'F' {
		FatalDepth(OutputCallDepth, string(p[DefaultPrefixLength:]))
	}
	return len(p), nil
}

// TestInfoWithSetOutputBySeverity helps ensure that we do not alter some of the assumptions
// used in the co-existence examples. Namely the depth of the stack where the klog v2 call
// originates from and the length of the prefix. So folks who use SetOutput are guaranteed
// what they see (depth = OutputCallDepth and prefix length = DefaultPrefixLength, see Write method above)
func TestInfoWithSetOutputBySeverity(t *testing.T) {
	setFlags()
	defer logging.swap(logging.newBuffers())
	LogToStderr(false)
	SetOutputBySeverity("INFO", klogWriter{})
	defer logging.swap(logging.newBuffers())
	Info("test")
	if !contains(infoLog, "I", t) {
		t.Errorf("Info has wrong character: [%q]", contents(infoLog))
	}
	if !contains(infoLog, "klog_overrides_test.go:47]", t) { // :XX matches the line number for Info method above
		t.Errorf("Info failed : got [%q]", contents(infoLog))
	}
}
