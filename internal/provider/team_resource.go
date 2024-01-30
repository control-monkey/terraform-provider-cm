package provider

import (
	"context"
	"fmt"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/team"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &TeamResource{}

func NewTeamResource() resource.Resource {
	return &TeamResource{}
}

type TeamResource struct {
	client *ControlMonkeyAPIClient
}

func (r *TeamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (r *TeamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys teams.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the team.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the team.",
				Required:            true,
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"custom_idp_id": schema.StringAttribute{
				MarkdownDescription: "Custom ID for identity provider (IdP)",
				Optional:            true,
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *TeamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get current state
	var state team.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.team.ReadTeam(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read team %s", id), fmt.Sprintf("%s", err))
		return
	}

	team.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *TeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan team.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := team.Converter(&plan, nil, commons.CreateConverter)

	res, err := r.client.Client.team.CreateTeam(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Team creation failed",
			fmt.Sprintf("failed to create team, error: %s", err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(controlmonkey.StringValue(res.ID))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan team.ResourceModel
	var state team.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body, _ := team.Converter(&plan, &state, commons.UpdateConverter)

	_, err := r.client.Client.team.UpdateTeam(ctx, id, body)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Team update failed",
			fmt.Sprintf("failed to update team %s, error: %s", id, err),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state team.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	_, err := r.client.Client.team.DeleteTeam(ctx, id)

	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Team deletion failed",
			fmt.Sprintf("Failed to delete team %s, error: %s", id, err),
		)
		return
	}
}

func (r *TeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
