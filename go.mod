module github.com/civo/bizaar-daemon

go 1.15

require (
	github.com/civo/bizaar-operator v0.0.0-20210209044545-f6c15afce6dc
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.3
)

// TODO remove this
// replace github.com/civo/bizaar-operator => /Users/zulh/civo/bizaar-operator
