package provider

import (
	"context"
	"fmt"
	"slices"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	tfCustomAbacConfiguration "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/custom_abac_configuration"
	cm_stringvalidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &CustomAbacConfigurationResource{}

func NewCustomAbacConfigurationResource() resource.Resource {
	return &CustomAbacConfigurationResource{}
}

type CustomAbacConfigurationResource struct {
	client *ControlMonkeyAPIClient
}

func (r *CustomAbacConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_abac_configuration"
}

func (r *CustomAbacConfigurationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys custom abac configurations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the custom abac configuration.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_abac_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the group in the IDP",
				Required:            true,
				Validators: []validator.String{
					cm_stringvalidators.NotBlank(),
					stringvalidator.LengthAtMost(255),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the custom abac configuration.",
				Required:            true,
				Validators: []validator.String{
					cm_stringvalidators.NotBlank(),
				},
			},
			"roles": schema.ListNestedAttribute{
				MarkdownDescription: "List of roles of the custom abac configuration.",
				Required:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"org_id": schema.StringAttribute{
							MarkdownDescription: "The Organization ID in ControlMonkey. It can be found [here](https://console.controlmonkey.io/app/organization/setting?id=idp). Example, `o-123`",
							Required:            true,
							Validators: []validator.String{
								cm_stringvalidators.NotBlank(),
							},
						},
						"org_role": schema.StringAttribute{
							MarkdownDescription: "The type of the role. Find supported types [here](https://docs.controlmonkey.io/controlmonkey-api/api-enumerations#custom-abac-org-role-types)",
							Required:            true,
							Validators: []validator.String{
								cm_stringvalidators.NotBlank(),
								stringvalidator.NoneOf("runner"),
							},
						},
						"team_ids": schema.ListAttribute{
							MarkdownDescription: "List of teams to assign the role to. This property cannot be used when `org_role` is set to admin/viewer",
							ElementType:         types.StringType,
							Optional:            true,
							Validators:          commons.ValidateUniqueNotEmptyListWithNoBlankValues(),
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *CustomAbacConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomAbacConfigurationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data tfCustomAbacConfiguration.ResourceModel

	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		return
	}

	roles := data.Roles

	if roles != nil {
		for _, r := range roles {
			orgRole := r.OrgRole

			if helpers.IsKnown(orgRole) {
				rolesWithoutAssignedTeams := []string{cmTypes.RoleAdmin, cmTypes.RoleViewer}

				if slices.Contains(rolesWithoutAssignedTeams, orgRole.ValueString()) {
					if r.TeamIds.IsNull() == false {
						resp.Diagnostics.AddError(
							validationError, fmt.Sprintf("Teams not allowed with %s/%s role", cmTypes.RoleAdmin, cmTypes.RoleViewer),
						)
					}
				}
			}
		}
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *CustomAbacConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state tfCustomAbacConfiguration.ResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.customAbacConfiguration.ReadCustomAbacConfiguration(ctx, id)

	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(resourceNotFoundError, fmt.Sprintf("Custom abac configuration '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read custom abac configuration '%s'", id), err.Error())
		return
	}

	tfCustomAbacConfiguration.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *CustomAbacConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan tfCustomAbacConfiguration.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := tfCustomAbacConfiguration.Converter(&plan, nil, commons.CreateConverter)

	res, err := r.client.Client.customAbacConfiguration.CreateCustomAbacConfiguration(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceCreationFailedError,
			fmt.Sprintf("failed to create custom abac configuration, error: %s", err.Error()),
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

func (r *CustomAbacConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan tfCustomAbacConfiguration.ResourceModel
	var state tfCustomAbacConfiguration.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body, _ := tfCustomAbacConfiguration.Converter(&plan, &state, commons.UpdateConverter)

	_, err := r.client.Client.customAbacConfiguration.UpdateCustomAbacConfiguration(ctx, id, body)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.Diagnostics.AddError(resourceNotFoundError, fmt.Sprintf("Custom abac configuration '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(
			resourceUpdateFailedError,
			fmt.Sprintf("failed to update custom abac configuration %s, error: %s", id, err),
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

func (r *CustomAbacConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state tfCustomAbacConfiguration.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	_, err := r.client.Client.customAbacConfiguration.DeleteCustomAbacConfiguration(ctx, id)

	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			resourceDeletionFailedError,
			fmt.Sprintf("Failed to delete custom abac configuration %s, error: %s", id, err),
		)
		return
	}
}

func (r *CustomAbacConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
