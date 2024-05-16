package provider

import (
	"context"
	"fmt"
	tfBlueprint "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/blueprint_data"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &BlueprintDataSource{}

func NewBlueprintDataSource() datasource.DataSource {
	return &BlueprintDataSource{}
}

type BlueprintDataSource struct {
	client *ControlMonkeyAPIClient
}

func (r *BlueprintDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blueprint"
}

func (r *BlueprintDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the blueprint.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("id"), path.MatchRoot("name")),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the blueprint.",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *BlueprintDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (r *BlueprintDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	//Get current state
	var state tfBlueprint.ResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	blueprintId := state.ID.ValueStringPointer()
	name := state.Name.ValueStringPointer()
	res, err := r.client.Client.blueprint.ListBlueprints(ctx, blueprintId, name)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read blueprint"), fmt.Sprintf("%s", err))
		return
	} else if len(res) == 0 {
		resp.Diagnostics.AddError(fmt.Sprintf(resourceNotFoundError), fmt.Sprintf(blueprintNotFoundError))
		return
	} else if len(res) > 1 {
		resp.Diagnostics.AddError(fmt.Sprintf(multipleEntitiesError), fmt.Sprintf("Found multiple blueprints with name '%s'; use additional constraints to reduce matches to a single match", *name))
		return
	}

	tfBlueprint.UpdateStateAfterRead(res[0], &state, &resp.Diagnostics)

	// Set refreshed state
	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
