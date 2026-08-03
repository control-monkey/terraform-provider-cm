package provider

import (
	"context"
	"fmt"

	sdkOrganizationPermissions "github.com/control-monkey/controlmonkey-sdk-go/services/organization_permissions"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	organizationPermission "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/organization_permissions"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OrganizationPermissionsResource{}

func NewOrganizationPermissionsResource() resource.Resource {
	return &OrganizationPermissionsResource{}
}

type OrganizationPermissionsResource struct {
	client *ControlMonkeyAPIClient
}

func (r *OrganizationPermissionsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_permissions"
}

func (r *OrganizationPermissionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates and destroys organization level permissions of a team." +
			" Each permission grants the team an organization role, which must be a `cm_custom_role` whose `type` is `organizationRole`." +
			" For more information: [ControlMonkey Documentation](https://docs.controlmonkey.io/administration/users-and-roles/custom-roles)\n\n" +
			"~> **Note** Organization permissions can only be managed for teams. The ControlMonkey API does not expose a way to read" +
			" the organization permissions of a user or of a programmatic user, so those cannot be managed by Terraform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the team.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"permissions": commons.WithNullSetDefault(schema.SetNestedAttribute{
				MarkdownDescription: "List of permissions",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"custom_role_id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of a custom role whose `type` is `organizationRole`.",
							Required:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
					},
				},
			}),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *OrganizationPermissionsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationPermissionsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data organizationPermission.ResourceModel

	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		return
	}

	if len(data.Permissions) > 0 {
		identifiers := interfaces.GetIdentifiers(data.Permissions)

		if helpers.IsUnique(identifiers) == false {
			duplicates := helpers.FindDuplicates(identifiers, false)
			for _, d := range duplicates {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("Custom role '%s' appears more than once", organizationPermission.CleanIdentifier(d)))
			}
		}
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *OrganizationPermissionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state organizationPermission.ResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.organizationPermissions.ListOrganizationPermissions(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(teamNotFoundError, fmt.Sprintf("Team '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read organization permissions for team '%s'", id), err.Error())
		return
	}

	organizationPermission.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *OrganizationPermissionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan organizationPermission.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := organizationPermission.Merge(&plan, nil, commons.CreateConverter)
	teamId := plan.TeamId

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, teamId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = teamId

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update handles permission changes. The API has no update endpoint, so a changed permission is
// a delete followed by a create.
func (r *OrganizationPermissionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan organizationPermission.ResourceModel
	var state organizationPermission.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := organizationPermission.Merge(&plan, &state, commons.UpdateMerger)
	teamId := plan.TeamId

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, teamId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, teamId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *OrganizationPermissionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state organizationPermission.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := organizationPermission.Merge(nil, &state, commons.DeleteMerger)
	teamId := state.TeamId

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, teamId.ValueString())
	resp.Diagnostics.Append(diags...)
}

func (r *OrganizationPermissionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *OrganizationPermissionsResource) createEntities(ctx context.Context, entitiesToCreate []*sdkOrganizationPermissions.OrganizationPermission, teamId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Granting %d organization permissions to team '%s'.", len(entitiesToCreate), teamId))

	for _, e := range entitiesToCreate {
		_, err := r.client.Client.organizationPermissions.CreateOrganizationPermission(ctx, e)

		if err != nil {
			customRoleId := *e.CustomRoleId
			if commons.IsAlreadyExistResponseError(err) {
				tflog.Info(ctx, fmt.Sprintf("Custom role '%s' is already granted to team '%s'. No operation was made.", customRoleId, teamId))
			} else if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to grant custom role '%s' to team '%s'. Error: %s", customRoleId, teamId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to grant custom role '%s' to team '%s'", customRoleId, teamId),
						err.Error()),
				}
			}
		}
	}

	return retVal
}

func (r *OrganizationPermissionsResource) deleteEntities(ctx context.Context, entitiesToDelete []*sdkOrganizationPermissions.OrganizationPermission, teamId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Revoking %d organization permissions from team '%s'.", len(entitiesToDelete), teamId))

	for _, e := range entitiesToDelete {
		_, err := r.client.Client.organizationPermissions.DeleteOrganizationPermission(ctx, e)

		if err != nil {
			customRoleId := *e.CustomRoleId
			if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to revoke custom role '%s' from team '%s'. Error: %s", customRoleId, teamId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to revoke custom role '%s' from team '%s'", customRoleId, teamId),
						err.Error()),
				}
			}
		}
	}

	return retVal
}
