package api

import (
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/terraform-provider-runtime/pkg/client"
	"github.com/crossplane/terraform-provider-runtime/pkg/plugin"
	"github.com/hashicorp/terraform/providers"
)

// Update syncs with an existing resource and modifies mutable values
func Update(p *client.Provider, inv *plugin.Invoker, res resource.Managed) (resource.Managed, error) {
	s, err := SchemaForInvoker(p, inv)
	if err != nil {
		return nil, err
	}
	encoded, err := inv.EncodeCty(res, s)
	if err != nil {
		return nil, err
	}
	tfName := inv.TerraformResourceName()

	prior, err := Read(p, inv, res)
	if err != nil {
		return nil, err
	}
	priorEncoded, err := inv.EncodeCty(prior, s)
	if err != nil {
		return nil, err
	}

	// TODO: research how/if the major providers are using Config
	// same goes for the private state blobs that are shuffled around
	req := providers.ApplyResourceChangeRequest{
		TypeName:   tfName,
		PriorState: priorEncoded,
		// TODO: For the purposes of Create, I am assuming that it's fine for
		// Config and PlannedState to be the same
		Config:       encoded,
		PlannedState: encoded,
	}
	resp := p.GRPCProvider.ApplyResourceChange(req)
	if resp.Diagnostics.HasErrors() {
		return res, resp.Diagnostics.NonFatalErr()
	}

	return inv.DecodeCty(res, resp.NewState, s)
}
