package client

import (
	"context"
	"io/ioutil"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/hashicorp/terraform/configs/configschema"
	tfplugin "github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/providers"
	"github.com/zclconf/go-cty/cty"
	kubeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

// FakeTerraformVersion is a nice lie we tell to providers to keep them happy
// TODO: is there a more sane way to negotiate version compat w/ providers?
const FakeTerraformVersion string = "v0.12.26"

// Provider wraps grpcProvider with some additional metadata like the provider name
type Provider struct {
	GRPCProvider *tfplugin.GRPCProvider
	Name         string
	Config       ProviderConfig
}

// ProviderConfig models the on-disk yaml config for providers
type ProviderConfig struct {
	TerraformConfig cty.Value `json:"config"`
}

// Configure calls the provider's grpc configuration interface,
// also translating the ProviderConfig structure to the
// Provider's encoded HCL representation.
func (p *Provider) Configure(cfg map[string]cty.Value) error {
	schema, err := GetProviderSchema(p)
	if err != nil {
		return err
	}
	// TODO: does not address nested blocks in the config
	for key, attr := range schema.Attributes {
		if _, ok := cfg[key]; !ok {
			switch attr.Type.FriendlyName() {
			case "string":
				cfg[key] = cty.NullVal(cty.String)
				continue
			case "bool":
				cfg[key] = cty.NullVal(cty.Bool)
				continue
			case "list of string":
				cfg[key] = cty.ListValEmpty(cty.String)
			default:
				cfg[key] = cty.NullVal(cty.EmptyObject)
			}
		}
	}
	ctyCfg := cty.ObjectVal(cfg)
	cfgReq := providers.ConfigureRequest{
		TerraformVersion: FakeTerraformVersion,
		Config:           ctyCfg,
	}
	cfgResp := p.GRPCProvider.Configure(cfgReq)
	if cfgResp.Diagnostics.HasErrors() {
		return cfgResp.Diagnostics.Err()
	}

	return nil
}

// ReadProviderConfigFile reads a yaml-formatted provider config and unmarshals
// it into a ProviderConfig, which knows how to generate the serialized
// provider config that a terraform provider expects.
func ReadProviderConfigFile(path string) (ProviderConfig, error) {
	cfg := ProviderConfig{}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = yaml.UnmarshalStrict(content, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, err
}

// NewProvider constructs a Provider, which is a container type, holding a
// terraform provider plugin grpc client, as well as metadata about this provider
// instance, eg its configuration and type.
func NewProvider(providerName string, pluginDir string, cfg map[string]cty.Value) (*Provider, error) {
	grpc, err := NewGRPCProvider(providerName, pluginDir)
	if err != nil {
		return nil, err
	}
	provider := &Provider{
		Name:         providerName,
		GRPCProvider: grpc,
	}
	err = provider.Configure(cfg)

	return provider, err
}

func GetProviderSchema(p *Provider) (*configschema.Block, error) {
	resp := p.GRPCProvider.GetSchema()
	if resp.Diagnostics.HasErrors() {
		return resp.Provider.Block, resp.Diagnostics.NonFatalErr()
	}
	return resp.Provider.Block, nil
}

type Initializer func(context.Context, resource.Managed, *RuntimeOptions, kubeclient.Client) (*Provider, error)

var DefaultProviderPoolSize = 5

type RuntimeOptions struct {
	PoolSize        int
	PluginDirectory string
}

var DefaultPluginDirectory string = "/Users/kasey/src/crossplane/provider-terraform-gcp/.terraform/plugins/darwin_amd64/"

func (ro *RuntimeOptions) GetPluginDirectory() string {
	return DefaultPluginDirectory
}

func (ro *RuntimeOptions) WithPluginDirectory(dir string) *RuntimeOptions {
	ro.PluginDirectory = dir
	return ro
}

func (ro *RuntimeOptions) WithPoolSize(size int) *RuntimeOptions {
	ro.PoolSize = size
	return ro
}

func NewRuntimeOptions() *RuntimeOptions {
	return &RuntimeOptions{}
}
