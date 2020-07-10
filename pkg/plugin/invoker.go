package plugin

import (
	"fmt"

	xpresource "github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/hashicorp/terraform/providers"
	"github.com/zclconf/go-cty/cty"
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

type Invoker struct {
	ft Implementation
}

func (a *Invoker) GVK() k8schema.GroupVersionKind {
	return a.ft.GVK
}

func (a *Invoker) TerraformResourceName() string {
	return a.ft.TerraformResourceName
}

func (a *Invoker) EncodeCty(r xpresource.Managed, s *providers.Schema) (cty.Value, error) {
	if a.ft.CtyEncoder == nil {
		return cty.Value{}, fmt.Errorf("Cannot lookup EncodeCty for GVK=%s", a.ft.GVK.String())
	}
	return a.ft.CtyEncoder.EncodeCty(r, s)
}

func (a *Invoker) DecodeCty(r xpresource.Managed, v cty.Value, s *providers.Schema) (xpresource.Managed, error) {
	if a.ft.CtyDecoder == nil {
		return nil, fmt.Errorf("Cannot lookup DecodeCty for GVK=%s", a.ft.GVK.String())
	}
	return a.ft.CtyDecoder.DecodeCty(r, v, s)
}

func (a *Invoker) SchemeBuilder() (*scheme.Builder, error) {
	if a.ft.SchemeBuilder == nil {
		return nil, fmt.Errorf("Cannot lookup SchemeBuilder for GVK=%s", a.ft.GVK.String())
	}
	return a.ft.SchemeBuilder, nil
}

func (a *Invoker) UnmarshalResourceYaml(b []byte) (xpresource.Managed, error) {
	if a.ft.ResourceYAMLUnmarshaller == nil {
		return nil, fmt.Errorf("Cannot lookup UnmarshalResourceCallback for GVK=%s", a.ft.GVK.String())
	}
	return a.ft.ResourceYAMLUnmarshaller.UnmarshalResourceYAML(b)
}

func (a *Invoker) MarshalResourceYaml(r xpresource.Managed) ([]byte, error) {
	if a.ft.ResourceYAMLMarshaller == nil {
		return nil, fmt.Errorf("Cannot lookup YamlEncodeCallback for GVK=%s", a.ft.GVK.String())
	}
	return a.ft.ResourceYAMLMarshaller.MarshalResourceYAML(r)
}

func (a *Invoker) MergeResources(f xpresource.Managed, t xpresource.Managed) (MergeDescription, error) {
	if a.ft.ResourceMerger == nil {
		return MergeDescription{}, fmt.Errorf("Cannot lookup MergeResources() for GVK=%s", a.ft.GVK.String())
	}
	return a.ft.ResourceMerger.MergeResources(f, t), nil
}
