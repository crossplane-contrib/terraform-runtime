package api

import (
	"github.com/crossplane-contrib/terraform-runtime/pkg/client"
	"github.com/crossplane-contrib/terraform-runtime/pkg/plugin"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/hashicorp/terraform/providers"
	"github.com/zclconf/go-cty/cty"
)

// Delete deletes the given resource from the provider
// In terraform slang this is expressed as asking the provider
// to act on a Nil planned state.
func Delete(p *client.Provider, inv *plugin.Invoker, res resource.Managed) error {
	s, err := SchemaForInvoker(p, inv)
	if err != nil {
		return err
	}
	encoded, err := inv.EncodeCty(res, s)
	if err != nil {
		return err
	}

	req := providers.ApplyResourceChangeRequest{
		TypeName:   inv.TerraformResourceName(),
		PriorState: encoded,
		// TODO: For the purposes of Delete, I am assuming that it's fine for
		// Config and PlannedState to be the same
		Config:       cty.NullVal(s.Block.ImpliedType()),
		PlannedState: cty.NullVal(s.Block.ImpliedType()),
		ProviderMeta: cty.NullVal(cty.DynamicPseudoType),
	}
	resp := p.GRPCProvider.ApplyResourceChange(req)
	if resp.Diagnostics.HasErrors() {
		return resp.Diagnostics.NonFatalErr()
	}
	return nil
}
