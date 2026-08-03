package provider

import (
	"context"
	"fmt"

	tfCustomFlow "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/custom_flow_data"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &CustomFlowDataSource{}

func NewCustomFlowDataSource() datasource.DataSource {
	return &CustomFlowDataSource{}
}

type CustomFlowDataSource struct {
	client *ControlMonkeyAPIClient
}

func (r *CustomFlowDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_flow"
}

func (r *CustomFlowDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the custom flow.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("id"), path.MatchRoot("name")),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the custom flow.",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *CustomFlowDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (r *CustomFlowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	//Get current state
	var state tfCustomFlow.ResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	customFlowId := state.ID.ValueStringPointer()
	name := state.Name.ValueStringPointer()
	res, err := r.client.Client.customFlow.ListCustomFlows(ctx, customFlowId, name, nil, nil)

	if err != nil {
		resp.Diagnostics.AddError("Failed to read custom flow", err.Error())
		return
	} else if len(res) == 0 {
		resp.Diagnostics.AddError(resourceNotFoundError, customFlowNotFoundError)
		return
	} else if len(res) > 1 {
		resp.Diagnostics.AddError(multipleEntitiesError,
			fmt.Sprintf("Found %d custom flows matching name '%s' and id '%s'. Use additional constraints to reduce matches to a single match.",
				len(res), state.Name.ValueString(), state.ID.ValueString()))
		return
	}

	tfCustomFlow.UpdateStateAfterRead(res[0], &state, &resp.Diagnostics)

	// Set refreshed state
	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
