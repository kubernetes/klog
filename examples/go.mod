module example

go 1.13

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	k8s.io/klog/v2 v2.30.0
)

replace k8s.io/klog/v2 => ../
