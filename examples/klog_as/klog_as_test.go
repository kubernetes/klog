// Copyright 2021 The Kubernetes Authors.
//
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

package klog_test

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-logr/logr"

	"k8s.io/klog/v2"
)

type myStruct struct {
	some, data string
}

func ExampleAs() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "false")
	flag.Set("skip_headers", "true")
	klog.SetOutput(os.Stdout)

	item := myStruct{
		some: "thing",
		data: "someone",
	}
	klog.InfoS("No formatting", "item", item)

	as := klog.As(
		func() string {
			return fmt.Sprintf("{Some: %q Data: %q}", item.some, item.data)
		},
		func() interface{} {
			return struct {
				Data, Some string
			}{
				Some: item.some,
				Data: item.data,
			}
		},
	)
	klog.InfoS("With stringer", "item", as)

	// We don't have a logger in klog which uses MarshalLog, but we can
	// test its behavior by invoking it directly.
	klog.InfoS("With marshaler", "item", as.(logr.Marshaler).MarshalLog())

	// Not a useful call, but klog.As tolerates it.
	klog.InfoS("Callbacks should never be nil", "item", klog.As(nil, nil))

	klog.InfoS("hello", "obj", klog.AsText(func() string { return "world" }))

	// Output:
	// "No formatting" item={some:thing data:someone}
	// "With stringer" item="{Some: \"thing\" Data: \"someone\"}"
	// "With marshaler" item={Data:someone Some:thing}
	// "Callbacks should never be nil" item=<nil>
	// "hello" obj="world"
}
