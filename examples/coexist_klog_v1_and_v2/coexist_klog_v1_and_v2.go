package main

import (
	"flag"

	klogv1 "k8s.io/klog"
	klogv2 "k8s.io/klog/v2"
)

// OutputCallDepth is the stack depth where we can find the origin of this call
const OutputCallDepth = 6
// log prefix that we have to strip out
const DefaultPrefixLength = 53

// klogWriter is used in SetOutputBySeverity call below to redirect
// any calls to klogv1 to end up in klogv2
type klogWriter struct{}

func (kw klogWriter) Write(p []byte) (n int, err error) {
	if len(p) < DefaultPrefixLength {
		klogv2.InfoDepth(OutputCallDepth, string(p))
		return len(p), nil
	}
	if p[0] == 'I' {
		klogv2.InfoDepth(OutputCallDepth, string(p[DefaultPrefixLength:]))
	} else if p[0] == 'W' {
		klogv2.WarningDepth(OutputCallDepth, string(p[DefaultPrefixLength:]))
	} else if p[0] == 'E' {
		klogv2.ErrorDepth(OutputCallDepth, string(p[DefaultPrefixLength:]))
	} else if p[0] == 'F' {
		klogv2.FatalDepth(OutputCallDepth, string(p[DefaultPrefixLength:]))
	} else {
		klogv2.InfoDepth(OutputCallDepth, string(p[DefaultPrefixLength:]))
	}
	return len(p) - DefaultPrefixLength, nil
}

func main() {
	// Assume that we have to switch klogv2 on from now on everywhere.
	klogv2.InitFlags(nil)

	// These following are just things you do to ensure that you stream
	// all the things to the log file instead of stdout/stderr
	flag.Set("logtostderr", "false")     // By default klog logs to stderr, switch that off
	flag.Set("alsologtostderr", "false") // false is default, but this is informative
	flag.Set("stderrthreshold", "FATAL") // stderrthreshold defaults to ERROR, we don't want anything in stderr
	flag.Set("log_file", "myfile.log")   // log to a file
	flag.Parse()

	// BEGIN : hack to redirect klogv1 calls to klog v2
	// Tell klog to log into STDERR. Otherwise, we risk
	// certain kinds of API errors getting logged into a directory not
	// available in a `FROM scratch` Docker container, causing us to abort
	var flags flag.FlagSet
	klogv1.InitFlags(&flags)
	flags.Set("logtostderr", "false")
	klogv1.SetOutputBySeverity("INFO", klogWriter{})
	// END : hack to redirect klogv1 calls to klog v2

	// Now you can mix klogv1 and v2 in the same code base
	klogv2.Info("hello from klog (v2)!")
	klogv1.Info("hello from klog (v1)!")
	klogv1.Warning("beware from klog (v1)!")
	klogv1.Error("error from klog (v1)!")
	klogv2.Info("nice to meet you (v2)")
	klogv2.Flush()
}
