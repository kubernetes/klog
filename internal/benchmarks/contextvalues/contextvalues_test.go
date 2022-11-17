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

package contextvalues

import (
	"context"
	"testing"

	"github.com/go-logr/logr"

	"k8s.io/klog/v2"
)

const iterationsPerOp = 100

// 1% of the Info calls are invoked, all of those call the LogSink.
func BenchmarkNewContext1Percent(b *testing.B) {
	ctx := klog.NewContext(context.Background(), discardWithV(1))
	defer func() {
		expected := int64(iterationsPerOp) / 100 * int64(b.N)
		if infoCalls != expected {
			b.Errorf("expected %d calls to Info, got %d", expected, infoCalls)
		}
		infoCalls = 0
	}()

	// Each iteration is expected to do exactly the same thing,
	// in particular do the same number of allocs.
	for i := 0; i < b.N; i++ {
		// Therefore we repeat newContext a certain number of
		// times. Individual repetitions are allowed to sometimes log
		// and sometimes not, but the overall execution is the same for
		// every outer loop iteration.
		for j := 0; j < iterationsPerOp; j++ {
			newContext(ctx, j, 100, 0)
		}
	}
}

// 100% of the Info calls are invoked, none of those call the LogSink.
func BenchmarkNewContext100PercentDisabled(b *testing.B) {
	ctx := klog.NewContext(context.Background(), discardWithV(1))
	defer func() {
		expected := int64(0)
		if infoCalls != expected {
			b.Errorf("expected %d calls to Info, got %d", expected, infoCalls)
		}
		infoCalls = 0
	}()

	for i := 0; i < b.N; i++ {
		for j := 0; j < iterationsPerOp; j++ {
			newContext(ctx, j, 1, 2)
		}
	}
}

// 100% of the Info calls are invoked, all of those call the LogSink.
func BenchmarkNewContext100Percent(b *testing.B) {
	ctx := klog.NewContext(context.Background(), discardWithV(1))
	defer func() {
		expected := int64(b.N) * iterationsPerOp
		if infoCalls != expected {
			b.Errorf("expected %d calls to Info, got %d", expected, infoCalls)
		}
		infoCalls = 0
	}()

	for i := 0; i < b.N; i++ {
		for j := 0; j < iterationsPerOp; j++ {
			newContext(ctx, j, 1, 0)
		}
	}
}

func BenchmarkNewContext1PercentValues(b *testing.B) {
	initContextValues(b)
	ctx := klog.NewContext(context.Background(), discardWithV(1))
	defer func() {
		expected := int64(iterationsPerOp) / 100 * int64(b.N)
		if infoCalls != expected {
			b.Errorf("expected %d calls to Info, got %d", expected, infoCalls)
		}
		infoCalls = 0
	}()

	for i := 0; i < b.N; i++ {
		for j := 0; j < iterationsPerOp; j++ {
			newContextValues(ctx, j, 100, 0)
		}
	}
}

func BenchmarkNewContext100PercentDisabledValues(b *testing.B) {
	initContextValues(b)
	ctx := klog.NewContext(context.Background(), discardWithV(1))
	defer func() {
		expected := int64(0)
		if infoCalls != expected {
			b.Errorf("expected %d calls to Info, got %d", expected, infoCalls)
		}
		infoCalls = 0
	}()

	for i := 0; i < b.N; i++ {
		for j := 0; j < iterationsPerOp; j++ {
			newContextValues(ctx, j, 1, 2)
		}
	}
}

func BenchmarkNewContext100PercentValues(b *testing.B) {
	initContextValues(b)
	ctx := klog.NewContext(context.Background(), discardWithV(1))
	defer func() {
		expected := int64(b.N) * iterationsPerOp
		if infoCalls != expected {
			b.Errorf("expected %d calls to Info, got %d", expected, infoCalls)
		}
		infoCalls = 0
	}()

	for i := 0; i < b.N; i++ {
		for j := 0; j < iterationsPerOp; j++ {
			newContextValues(ctx, j, 1, 0)
		}
	}
}

type contextKey1 struct{}
type contextKey2 struct{}

func initContextValues(b *testing.B) {
	state := klog.CaptureState()
	b.Cleanup(state.Restore)
	klog.SetFromContextKeys(
		klog.ContextKey{contextKey1{}, "i"},
		klog.ContextKey{contextKey2{}, "j"},
	)
}

func newContext(ctx context.Context, j, mod, v int) {
	// This is the currently recommended way of adding a value to a context
	// and ensuring that all future log calls include it.  Trace IDs might
	// get handled like this.
	logger := klog.FromContext(ctx)
	logger = klog.LoggerWithValues(logger, "i", 1, "j", 2)
	ctx = context.WithValue(ctx, contextKey1{}, 1)
	ctx = context.WithValue(ctx, contextKey2{}, 2)
	ctx = klog.NewContext(ctx, logger)
	useContext(ctx, j, mod, v)
}

func newContextValues(ctx context.Context, j, mod, v int) {
	// This variant only adds to the context. It relies on
	// klog.FromContextKeys to log them.
	ctx = context.WithValue(ctx, contextKey1{}, 1)
	ctx = context.WithValue(ctx, contextKey2{}, 2)
	useContext(ctx, j, mod, v)
}

func useContext(ctx context.Context, j, mod, v int) {
	if j%mod == 0 {
		logger := klog.FromContext(ctx)
		logger.V(v).Info("ping", "string", "hello world", "int", 1, "float", 1.0)
	}
}

func discardWithV(v int) logr.Logger {
	logger := logr.New(&discardLogSink{v: v})
	return logger
}

var infoCalls int64

// discardLogSink is a LogSink that discards all messages but has
// a verbosity threshold to compare calls which call Info vs. those
// that don't.
type discardLogSink struct {
	v int
}

func (l discardLogSink) Init(logr.RuntimeInfo) {
}

func (l discardLogSink) Enabled(v int) bool {
	return v <= l.v
}

func (l discardLogSink) Info(int, string, ...interface{}) {
	infoCalls++
}

func (l discardLogSink) Error(error, string, ...interface{}) {
}

func (l discardLogSink) WithValues(...interface{}) logr.LogSink {
	return l
}

func (l discardLogSink) WithName(string) logr.LogSink {
	return l
}
