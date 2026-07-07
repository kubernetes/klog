/*
Copyright 2024 The Kubernetes Authors.

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
	"reflect"
	"testing"

	"github.com/go-logr/logr"

	"k8s.io/klog/v2/internal/test"
)

// TestKObjSlice covers the two code paths that go through kobjSlice.process:
// the fmt.Stringer output used by text loggers ("String" sub-test) and the
// structured output used by logr backends ("MarshalLog" sub-test). Both must
// render a nil pointer entry as null rather than an empty ObjectRef.
func TestKObjSlice(t *testing.T) {
	pod1 := ObjectRef{Name: "pod-1", Namespace: "kube-system"}
	pod2 := ObjectRef{Name: "pod-2", Namespace: "kube-system"}
	for name, tc := range map[string]struct {
		arg         interface{}
		wantString  string
		wantMarshal interface{}
	}{
		"nil": {
			arg:         nil,
			wantString:  "[]",
			wantMarshal: []interface{}(nil),
		},
		"empty": {
			arg:         []test.KMetadataMock{},
			wantString:  "[]",
			wantMarshal: []interface{}{},
		},
		"objects": {
			arg: []test.KMetadataMock{
				{Name: "pod-1", NS: "kube-system"},
				{Name: "pod-2", NS: "kube-system"},
			},
			wantString:  "[kube-system/pod-1 kube-system/pod-2]",
			wantMarshal: []interface{}{pod1, pod2},
		},
		"untyped-nil-entry": {
			arg: []interface{}{
				test.KMetadataMock{Name: "pod-1", NS: "kube-system"},
				nil,
			},
			wantString:  "[kube-system/pod-1 <nil>]",
			wantMarshal: []interface{}{pod1, nil},
		},
		"typed-nil-pointer-entry": {
			arg: []*test.PtrKMetadataMock{
				{Name: "pod-1", NS: "kube-system"},
				nil,
			},
			wantString:  "[kube-system/pod-1 <nil>]",
			wantMarshal: []interface{}{pod1, nil},
		},
		"not-a-slice": {
			arg:         1,
			wantString:  "<KObjSlice needs a slice, got type int>",
			wantMarshal: "<KObjSlice needs a slice, got type int>",
		},
		"wrong-element-type": {
			arg:         []int{1, 2},
			wantString:  "<KObjSlice needs a slice of values implementing KMetadata, got type int>",
			wantMarshal: "<KObjSlice needs a slice of values implementing KMetadata, got type int>",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("String", func(t *testing.T) {
				got := KObjSlice(tc.arg).(fmt.Stringer).String()
				if got != tc.wantString {
					t.Errorf("KObjSlice(%#v).String() = %q, want %q", tc.arg, got, tc.wantString)
				}
			})
			t.Run("MarshalLog", func(t *testing.T) {
				got := KObjSlice(tc.arg).(logr.Marshaler).MarshalLog()
				if !reflect.DeepEqual(got, tc.wantMarshal) {
					t.Errorf("KObjSlice(%#v).MarshalLog() = %#v, want %#v", tc.arg, got, tc.wantMarshal)
				}
			})
		})
	}
}
