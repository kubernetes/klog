package main

import (
	"flag"
	"github.com/go-logr/logr"

	"k8s.io/klog"
	klogv2 "k8s.io/klog/v2"
)

type klogLogger struct{}

func (kw klogLogger) Info(msg string, keysAndValues ...interface{}) {
	// if you start using structured logging, please serialize as well
	klog.InfoDepth(klogv2.LoggerCallDepth, msg)
}

func (kw klogLogger) Enabled() bool {
	return true
}

func (kw klogLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	// if you start using structured logging, please serialize as well
	klog.ErrorDepth(klogv2.LoggerCallDepth, msg)
}

func (kw klogLogger) V(level int) logr.InfoLogger {
	// tweak the return value based in the level as needed
	return kw
}

func (kw klogLogger) WithValues(keysAndValues ...interface{}) logr.Logger {
	// if you start using structured logging, please serialize as well
	return kw
}

func (kw klogLogger) WithName(name string) logr.Logger {
	// if you start using structured logging, please serialize as well
	return kw
}

func main() {
	klog.InitFlags(nil)
	klogv2.LogToStderr(false)
	klogv2.SetLogger(klogLogger{})

	flag.Set("logtostderr", "false")     // By default klog logs to stderr, switch that off
	flag.Set("alsologtostderr", "false") // false is default, but this is informative
	flag.Set("stderrthreshold", "FATAL") // stderrthreshold defaults to ERROR, we don't want anything in stderr
	flag.Set("log_file", "myfile.log")   // log to a file
	flag.Parse()

	klog.Info("hello from klog v1!")
	klogv2.Info("hello from klog v2!")
	klogv2.Warning("beware from klog v2!")
	klogv2.Error("error from klog v2!")
	klog.Info("nice to meet you")
	klog.Flush()
}
