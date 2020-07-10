package api

import (
	"fmt"

	"github.com/crossplane/terraform-provider-runtime/pkg/client"
	"github.com/crossplane/terraform-provider-runtime/pkg/plugin"
	"github.com/hashicorp/terraform/providers"
	"github.com/pkg/errors"
)

func GetSchema(p *client.Provider) (map[string]providers.Schema, error) {
	resp := p.GRPCProvider.GetSchema()
	if resp.Diagnostics.HasErrors() {
		return nil, resp.Diagnostics.NonFatalErr()
	}

	return resp.ResourceTypes, nil
}

func SchemaForInvoker(p *client.Provider, inv *plugin.Invoker) (*providers.Schema, error) {
	schema, err := GetSchema(p)
	if err != nil {
		msg := "Failed to retrieve schema from provider in api.Read"
		return nil, errors.Wrap(err, msg)
	}
	tfName := inv.TerraformResourceName()
	s, ok := schema[tfName]
	if !ok {
		return nil, fmt.Errorf("Could not look up schema using terraform resource name=%s (for gvk=%s", tfName, inv.GVK().String())
	}

	return &s, nil
}
