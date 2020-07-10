package plugin

import (
	xpresource "github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/hashicorp/terraform/providers"
	"github.com/zclconf/go-cty/cty"
)

type ResourceYAMLMarshaller interface {
	MarshalResourceYAML(xpresource.Managed) ([]byte, error)
}

type ResourceYAMLUnmarshaller interface {
	UnmarshalResourceYAML([]byte) (xpresource.Managed, error)
}

type CtyEncoder interface {
	EncodeCty(xpresource.Managed, *providers.Schema) (cty.Value, error)
}

type CtyDecoder interface {
	DecodeCty(xpresource.Managed, cty.Value, *providers.Schema) (xpresource.Managed, error)
}
