module github.com/crossplane/terraform-provider-runtime

go 1.13

require (
	github.com/crossplane/crossplane v0.11.1
	github.com/crossplane/crossplane-runtime v0.9.0
	github.com/crossplane/crossplane-tools v0.0.0-20200412230150-efd0edd4565b
	github.com/hashicorp/go-plugin v1.3.0
	github.com/hashicorp/terraform v0.12.29
	github.com/pkg/errors v0.9.1
	github.com/zclconf/go-cty v1.5.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/controller-tools v0.2.4
	sigs.k8s.io/yaml v1.2.0
)
