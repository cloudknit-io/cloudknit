module github.com/compuzest/environment-operator

go 1.13

require (
	github.com/argoproj/argo v2.5.2+incompatible
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.2
)

replace sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.2.8

