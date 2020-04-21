package main

import (
	"flag"

	"k8s.io/klog"
	klogv2 "k8s.io/klog/v2"
)

type klogWriter struct{}

func (kw klogWriter) Write(p []byte) (n int, err error) {
	if p[0] == 'I' {
		klog.InfoDepth(klogv2.OutputCallDepth, string(p[klogv2.DefaultPrefixLength:]))
	} else if p[0] == 'W' {
		klog.WarningDepth(klogv2.OutputCallDepth, string(p[klogv2.DefaultPrefixLength:]))
	} else if p[0] == 'E' {
		klog.ErrorDepth(klogv2.OutputCallDepth, string(p[klogv2.DefaultPrefixLength:]))
	} else if p[0] == 'F' {
		klog.FatalDepth(klogv2.OutputCallDepth, string(p[klogv2.DefaultPrefixLength:]))
	}
	return len(p), nil
}

func main() {
	klog.InitFlags(nil)
	klogv2.LogToStderr(false)
	klogv2.SetOutputBySeverity("INFO", klogWriter{})

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
