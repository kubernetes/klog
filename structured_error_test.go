/*
Copyright The Kubernetes Authors.

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
	"errors"
	"reflect"
	"testing"
)

func TestErrorDetails(t *testing.T) {
	base := errors.New("base")

	for name, tc := range map[string]struct {
		err                error
		expectErrorString  string
		expectErrorDetails any
	}{
		"simple": {ErrorWithDetails(base, 42), "base", 42},
		"pair":   {ErrorWithDetails(ErrorWithDetails(base, "hello"), "world"), "base", []any{"hello", "world"}},
		"nested": {ErrorWithDetails(ErrorWithDetails(ErrorWithDetails(base, "hello"), "world"), "thanks"), "base", []any{"hello", "world", "thanks"}},

		"simple-func": {ErrorWithDetailsFunc(base, func() any { return 42 }), "base", 42},
		"pair-func":   {ErrorWithDetailsFunc(ErrorWithDetails(base, "hello"), func() any { return "world" }), "base", []any{"hello", "world"}},
		"nested-func": {ErrorWithDetailsFunc(ErrorWithDetails(ErrorWithDetails(base, "hello"), "world"), func() any { return "thanks" }), "base", []any{"hello", "world", "thanks"}},
	} {
		t.Run(name, func(t *testing.T) {
			if actual, expect := tc.err.Error(), tc.expectErrorString; actual != expect {
				t.Errorf("expected error string %q, got %q", expect, actual)
			}
			if actual, expect := tc.err.(ErrorDetailer).ErrorDetails(), tc.expectErrorDetails; !reflect.DeepEqual(actual, expect) {
				t.Errorf("expected error details %#v, got %#v", expect, actual)
			}
		})
	}
}
