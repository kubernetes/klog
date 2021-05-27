/*
Copyright 2021 The Kubernetes Authors.

SPDX-License-Identifier: Apache-2.0
*/

package example

import (
	"fmt"
	"testing"

	"github.com/go-logr/logr"
	logrtesting "github.com/go-logr/logr/testing"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/testinglogger"
)

func TestKlog(t *testing.T) {
	log := testinglogger.New(t)
	exampleOutput(log)
}

func TestLogr(t *testing.T) {
	log := logrtesting.NewTestLoggerWithOptions(t,
		logrtesting.Options{
			LogTimestamp: true,
			Verbosity:    5,
		},
	)
	exampleOutput(log)
}

type pair struct {
	a, b int
}

func (p pair) String() string {
	return fmt.Sprintf("(%d, %d)", p.a, p.b)
}

var _ fmt.Stringer = pair{}

type err struct {
	msg string
}

func (e err) Error() string {
	return "failed: " + e.msg
}

var _ error = err{}

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

func exampleOutput(log logr.Logger) {
	log.Info("hello world")
	log.Error(err{msg: "some error"}, "failed")
	log.V(1).Info("verbosity 1")
	log.WithName("main").WithName("helper").Info("with prefix")
	log.Info("key/value pairs",
		"int", 1,
		"float", 2.0,
		"pair", pair{a: 1, b: 2},
		"kobj", klog.KObj(kmeta{name: "sally", namespace: "kube-system"}),
	)
}
