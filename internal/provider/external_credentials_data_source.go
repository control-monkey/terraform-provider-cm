package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	tfExternalCredentials "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/external_credentials_data"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &ExternalCredentialsDataSource{}

func NewExternalCredentialsDataSource() datasource.DataSource {
	return &ExternalCredentialsDataSource{}
}

type ExternalCredentialsDataSource struct {
	client *ControlMonkeyAPIClient
}

func (r *ExternalCredentialsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_external_credentials"
}

func (r *ExternalCredentialsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"vendor": schema.StringAttribute{
				MarkdownDescription: "The external credentials vendor type (aws/azure/gcp/etc). Find supported vendors [here] (https://docs.controlmonkey.io/controlmonkey-api/api-enumerations#external-credentials-vendor-types).",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The Unique Id of the external credentials.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("id"), path.MatchRoot("name")),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The Name of the external credentials.",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *ExternalCredentialsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (r *ExternalCredentialsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	//Get current state
	var state tfExternalCredentials.ResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vendor := state.Vendor.ValueString()
	id := state.ID.ValueStringPointer()
	name := state.Name.ValueStringPointer()
	res, err := r.client.Client.externalCredentials.ListExternalCredentials(ctx, vendor, id, name)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read external credentials"), fmt.Sprintf("%s", err))
		return
	} else if len(res) == 0 {
		resp.Diagnostics.AddError(fmt.Sprintf(resourceNotFoundError), fmt.Sprintf(externalCredentialsNotFoundError))
		return
	} else if len(res) > 1 {
		resp.Diagnostics.AddError(fmt.Sprintf("Found multiple entities"), fmt.Sprintf("Found multiple external credentials with name '%s', and id '%s' use additional constraints to reduce matches to a single match", *name, *id))
		return
	}

	tfExternalCredentials.UpdateStateAfterRead(res[0], &state, &resp.Diagnostics)

	// Set refreshed state
	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
