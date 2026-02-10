package provider

import (
	"context"
	"fmt"

	tfExternalCredential "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/external_credential_data"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &ExternalCredentialDataSource{}

func NewExternalCredentialDataSource() datasource.DataSource {
	return &ExternalCredentialDataSource{}
}

type ExternalCredentialDataSource struct {
	client *ControlMonkeyAPIClient
}

func (r *ExternalCredentialDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_external_credential"
}

func (r *ExternalCredentialDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Unique Id of the external credential.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The Name of the external credential.",
				Required:            true,
			},
			"vendor": schema.StringAttribute{
				MarkdownDescription: "The vendor of the external credential (aws/azure/gcp/datadog/etc).",
				Required:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *ExternalCredentialDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ControlMonkeyAPIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ControlMonkeyAPIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Read refreshes the Terraform state with the latest data.
func (r *ExternalCredentialDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	//Get current state
	var state tfExternalCredential.ResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vendor := state.Vendor.ValueString()
	name := state.Name.ValueString()
	res, err := r.client.Client.externalCredential.ListExternalCredentials(ctx, vendor, name)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read external credential"), fmt.Sprintf("%s", err))
		return
	} else if len(res) == 0 {
		resp.Diagnostics.AddError(fmt.Sprintf(resourceNotFoundError), fmt.Sprintf(externalCredentialNotFoundError))
		return
	} else if len(res) > 1 {
		resp.Diagnostics.AddError(fmt.Sprintf("Found multiple entities"), fmt.Sprintf("Found multiple external credentials with name '%s'; use additional constraints to reduce matches to a single match", name))
		return
	}

	tfExternalCredential.UpdateStateAfterRead(res[0], &state, &resp.Diagnostics)

	// Set refreshed state
	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
