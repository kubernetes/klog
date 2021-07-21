package calldepth

import (
	"k8s.io/klog/v2"
)

// Putting these functions into a separate file makes it possible to validate that
// their source code file is *not* logged because of WithCallDepth(1).

func myInfo(l klog.Logger, msg string) {
	klog.WithCallDepth(l, 1).Info(msg)
}

func myInfo2(l klog.Logger, msg string) {
	myInfo(klog.WithCallDepth(l, 1), msg)
}
