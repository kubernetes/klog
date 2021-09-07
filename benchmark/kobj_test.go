/*
Copyright 2021 The Kubernetes Authors.

SPDX-License-Identifier: Apache-2.0
*/

package example

import (
	"testing"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

type kmeta struct {
	name, namespace string
}

func (k kmeta) GetName() string {
	return k.name
}

func (k kmeta) GetNamespace() string {
	return k.namespace
}

var _ klog.KMetadata = kmeta{}

// var obj = kmeta{name: "some-fake-name", namespace: "kube-system"}
var obj = v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "some-fake-name", Namespace: "kube-system"}}

var result string

func BenchmarkKObjByPointer(b *testing.B) {
	var s string
	for n := 0; n < b.N; n++ {
		s = klog.KObj(&obj).String()
	}
	result = s
}

func BenchmarkKObj2ByPointer(b *testing.B) {
	var s string
	for n := 0; n < b.N; n++ {
		s = klog.KObj2(&obj).String()
	}
	result = s
}

func BenchmarkSkipObjByPointer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		klog.V(10).InfoS("skipped", "obj", klog.KObj(&obj))
	}
}

func BenchmarkSkipObj2ByPointer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		klog.V(10).InfoS("skipped", "obj", klog.KObj2(&obj))
	}
}

func BenchmarkDiscardObjByPointer(b *testing.B) {
	log := logr.Discard()
	for n := 0; n < b.N; n++ {
		log.Info("skipped", klog.KObj(&obj))
	}
}

func BenchmarkDiscardObj2ByPointer(b *testing.B) {
	log := logr.Discard()
	for n := 0; n < b.N; n++ {
		log.Info("skipped", klog.KObj2(&obj))
	}
}
