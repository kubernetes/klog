module k8s.io/klog/examples/coexist_klog_v1

go 1.13

require (
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.0.0
)

replace k8s.io/klog/v2 => ../.. // Do not copy this line into your project, we use this for testing
