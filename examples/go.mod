module example

go 1.13

require (
	github.com/go-logr/logr v1.2.3
	github.com/go-logr/zapr v1.2.3
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	go.uber.org/goleak v1.1.12
	go.uber.org/zap v1.19.0
	k8s.io/klog/v2 v2.30.0
)

replace k8s.io/klog/v2 => ../

replace github.com/go-logr/logr => github.com/pohly/logr v1.0.1-0.20221206165918-68f59133d07f
