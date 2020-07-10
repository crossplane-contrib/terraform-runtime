package plugin

import (
	"github.com/crossplane/terraform-provider-runtime/pkg/client"
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

type ProviderInit struct {
	GVK           k8schema.GroupVersionKind
	SchemeBuilder *scheme.Builder
	Initializer   client.Initializer
}
