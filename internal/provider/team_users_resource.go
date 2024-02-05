package provider

import (
	"context"
	"fmt"
	"github.com/control-monkey/controlmonkey-sdk-go/services/team"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	teamUsers "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/team_users"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const teamNotFoundError = "Team not found"

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &TeamUsersResource{}

func NewTeamUsersResource() resource.Resource {
	return &TeamUsersResource{}
}

type TeamUsersResource struct {
	client *ControlMonkeyAPIClient
}

func (r *TeamUsersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_users"
}

func (r *TeamUsersResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys team users.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of this resource.",
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
			"users": schema.SetNestedAttribute{
				MarkdownDescription: "List of users in this team",
				Optional:            true,
				Computed:            true,
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{},
						},
						[]attr.Value{
							types.ObjectValueMust(
								map[string]attr.Type{}, map[string]attr.Value{}),
						},
					),
				),
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"email": schema.StringAttribute{
							MarkdownDescription: "Email of user.",
							Required:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *TeamUsersResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TeamUsersResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data teamUsers.ResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if len(data.Users) > 0 {
		filterOutUnknowns := func(u *teamUsers.UserModel) bool {
			return u.GetBlockIdentifier() != ""
		}
		knownUsers := helpers.Filter(data.Users, filterOutUnknowns)

		mapToIdentifier := func(u *teamUsers.UserModel) string {
			return u.GetBlockIdentifier()
		}
		identifiers := helpers.Map(knownUsers, mapToIdentifier)

		if helpers.IsUnique(identifiers) == false {
			duplicates := helpers.FindDuplicates(identifiers, true)
			resp.Diagnostics.AddError("User appears more than once", fmt.Sprintf("User with email '%s' appears more than once", duplicates[0]))
		}
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *TeamUsersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get current state
	var state teamUsers.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.team.ListTeamUsers(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.Diagnostics.AddWarning("Team not found", fmt.Sprintf("Team '%s' not found", id))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read teamUsers %s", id), fmt.Sprintf("%s", err))
		return
	}

	teamUsers.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *TeamUsersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan teamUsers.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := teamUsers.Merge(&plan, nil, commons.CreateMerger)
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

func (r *TeamUsersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan teamUsers.ResourceModel
	var state teamUsers.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := teamUsers.Merge(&plan, &state, commons.UpdateMerger)

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, plan.TeamId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, plan.TeamId.ValueString())
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

func (r *TeamUsersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state teamUsers.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := teamUsers.Merge(nil, &state, commons.DeleteMerger)

	if diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, state.TeamId.ValueString()); diags != nil {
		if diags.ErrorsCount() > 0 && diags[0].Detail() == teamNotFoundError {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.Append(diags...)
	}
}

func (r *TeamUsersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

//region Private Methods

func (r *TeamUsersResource) createEntities(ctx context.Context, entitiesToCreate []*team.TeamUser, teamId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Adding %d users to team '%s'.", len(entitiesToCreate), teamId))

	for _, e := range entitiesToCreate {
		_, err := r.client.Client.team.CreateTeamUser(ctx, e)

		if err != nil {
			if commons.IsAlreadyExistResponseError(err) {
				tflog.Info(ctx, fmt.Sprintf("User '%s' is already in team '%s'. No operation was made.", *e.UserEmail, *e.TeamId))
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Team user creation failed",
						fmt.Sprintf("failed to create team user, error: %s", err)),
				}
			}
		}
	}

	return retVal
}

func (r *TeamUsersResource) deleteEntities(ctx context.Context, entitiesToDelete []*team.TeamUser, teamId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Removing %d users from team '%s'.", len(entitiesToDelete), teamId))

	for _, e := range entitiesToDelete {
		_, err := r.client.Client.team.DeleteTeamUser(ctx, e)

		if err != nil {
			if commons.IsNotFoundResponseError(err) {
				if commons.DoesErrorContains(err, "Team not found") {
					return diag.Diagnostics{
						diag.NewErrorDiagnostic("Team not found", fmt.Sprintf("Team '%s' not found", *e.TeamId)),
					}
				} else if commons.DoesErrorContains(err, "User does not exist") {
					tflog.Info(ctx, fmt.Sprintf("User '%s' does not exists. Removing them from team '%s'.", *e.UserEmail, *e.TeamId))
				} else {
					tflog.Info(ctx, fmt.Sprintf("User '%s' already not in team '%s'.", *e.UserEmail, *e.TeamId))
				}
			}
		}
	}

	return retVal
}

//endregion
