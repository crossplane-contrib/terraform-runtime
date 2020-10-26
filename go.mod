module github.com/crossplane/terraform-provider-runtime

go 1.14

require (
	github.com/crossplane/crossplane v0.13.0
	github.com/crossplane/crossplane-runtime v0.10.0
	github.com/hashicorp/go-plugin v1.3.0
	github.com/hashicorp/terraform v0.13.5
	github.com/pkg/errors v0.9.1
	github.com/zclconf/go-cty v1.7.0
	k8s.io/apimachinery v0.18.8
	sigs.k8s.io/controller-runtime v0.6.2
	sigs.k8s.io/yaml v1.2.0
)

replace (
	k8s.io/api => k8s.io/api v0.18.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.6
	k8s.io/client-go => k8s.io/client-go v0.18.6
)
