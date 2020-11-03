module github.com/compuzest/environment-operator

go 1.13

require (
	github.com/argoproj/argo v2.5.2+incompatible
	github.com/argoproj/pkg v0.3.0 // indirect
	github.com/colinmarc/hdfs v1.1.3 // indirect
	github.com/go-logr/logr v0.2.0
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/robfig/cron v1.2.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	gopkg.in/jcmturner/aescts.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/dnsutils.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/gokrb5.v5 v5.3.0 // indirect
	gopkg.in/jcmturner/rpc.v0 v0.0.2 // indirect
	k8s.io/apimachinery v0.19.3
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.3
)

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v0.1.0
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1
	k8s.io/api => k8s.io/api v0.17.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.2
	k8s.io/client-go => k8s.io/client-go v0.17.2
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.5.1
	github.com/colinmarc/hdfs => github.com/colinmarc/hdfs v1.1.4-0.20180805212432-9746310a4d31
)
