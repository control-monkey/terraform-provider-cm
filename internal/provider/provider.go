package provider

import (
	"context"
	"github.com/control-monkey/terraform-provider-cm/version"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey/credentials"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey/featureflag"
)

// Ensure ControlMonkeyProvider satisfies various provider interfaces.
var _ provider.Provider = &ControlMonkeyProvider{}

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &ControlMonkeyProvider{}
}

// ControlMonkeyProvider defines the provider implementation.
type ControlMonkeyProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type ControlMonkeyAPIClient struct {
	Client *Client
}

// ControlMonkeyProviderModel describes the provider data model.
type ControlMonkeyProviderModel struct {
	Token types.String `tfsdk:"token"`
}

func (p *ControlMonkeyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cm"
	resp.Version = p.version
}

func (p *ControlMonkeyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				MarkdownDescription: "A programmatic user token for ControlMonkey. This can also be set via the `CONTROL_MONKEY_TOKEN` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *ControlMonkeyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ControlMonkeyProviderModel

	// Check environment variables
	token := os.Getenv(credentials.EnvCredentialsVarToken)

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check configuration data, which should take precedence over
	// environment variable data, if found.
	if data.Token.ValueString() != "" {
		token = data.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing CONTROL_MONKEY_TOKEN environment variable",
			"While configuring the provider, the environment variable CONTROL_MONKEY_TOKEN was not found.",
		)
		// Not returning early allows the logic to collect all errors.
	}

	config := Config{
		Token:            token,
		FeatureFlags:     os.Getenv(featureflag.EnvVar),
		terraformVersion: version.Version,
	}

	client, err := config.Client()
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Cannot create client with a given token",
			"Failed to create a client session with a given token",
		)
	}

	apiClient := &ControlMonkeyAPIClient{
		Client: client,
	}

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *ControlMonkeyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVariableResource,
		NewStackResource,
		NewNamespaceResource,
		NewTemplateResource,
		NewControlPolicyMappingResource,
		NewControlPolicyGroupMappingResource,
		NewTeamResource,
		NewTeamUsersResource,
		NewNamespacePermissionsResource,
		NewTemplateNamespaceMappingsResource,
		NewBlueprintNamespaceMappingsResource,
	}
}

func (p *ControlMonkeyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
