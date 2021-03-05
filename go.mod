module github.com/kubemart/kubemart-daemon

go 1.16

require (
	github.com/kubemart/kubemart-operator v0.0.43
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.3
)

// replace github.com/kubemart/kubemart-operator => /Users/zulh/kubemart/kubemart-operator
