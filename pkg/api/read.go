package api

import (
	"github.com/crossplane-contrib/terraform-runtime/pkg/client"
	"github.com/crossplane-contrib/terraform-runtime/pkg/plugin"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/hashicorp/terraform/providers"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

var ErrNotFound = errors.New("Resource not found")

// Read returns an up-to-date version of the resource
// TODO: If `id` is unset for a new resource, how do we figure out
// what value needs to be used as the id?
func Read(p *client.Provider, inv *plugin.Invoker, res resource.Managed) (resource.Managed, error) {
	s, err := SchemaForInvoker(p, inv)
	if err != nil {
		return nil, err
	}
	encoded, err := inv.EncodeCty(res, s)
	req := providers.ReadResourceRequest{
		TypeName:     inv.TerraformResourceName(),
		PriorState:   encoded,
		Private:      nil,
		ProviderMeta: cty.NullVal(cty.DynamicPseudoType),
	}
	resp := p.GRPCProvider.ReadResource(req)
	if resp.Diagnostics.HasErrors() {
		return res, resp.Diagnostics.NonFatalErr()
	}
	// should we persist resp.Private in a blob in the resource to use on the next call?
	// Risky since size is unbounded, but we might be matching core behavior more carefully
	if resp.NewState.IsNull() {
		return nil, ErrNotFound
	}
	return inv.DecodeCty(res, resp.NewState, s)
}
