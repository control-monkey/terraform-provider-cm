package provider

import (
	"context"
	"fmt"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	sdkNamespacePermissions "github.com/control-monkey/controlmonkey-sdk-go/services/namespace_permissions"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	tfNamespacePermissions "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/namespace_permissions"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	"strings"
)

const namespaceNotFoundError = "Namespace not found"

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &NamespacePermissionsResource{}

func NewNamespacePermissionsResource() resource.Resource {
	return &NamespacePermissionsResource{}
}

type NamespacePermissionsResource struct {
	client *ControlMonkeyAPIClient
}

func (r *NamespacePermissionsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_permissions"
}

func (r *NamespacePermissionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys namespace permissions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of this resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"namespace_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the namespace.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"permissions": schema.SetNestedAttribute{
				MarkdownDescription: "Specifies a list of permissions granted to this namespace.",
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
						"user_email": schema.StringAttribute{
							MarkdownDescription: "Email address of the user.",
							Optional:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
						"programmatic_username": schema.StringAttribute{
							MarkdownDescription: "Username of the programmatic user.",
							Optional:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
						"team_id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the team.",
							Optional:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
						"role": schema.StringAttribute{
							MarkdownDescription: fmt.Sprintf("The role that is associated with this permission. Allowed values: %s.", helpers.EnumForDocs(cmTypes.NamespaceRoleTypes)),
							Optional:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
								stringvalidator.OneOf(cmTypes.NamespaceRoleTypes...),
							},
						},
						"custom_role_id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the custom role.",
							Optional:            true,
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
func (r *NamespacePermissionsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NamespacePermissionsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data tfNamespacePermissions.ResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if len(data.Permissions) > 0 {
		//Validate constraints
		for _, permission := range data.Permissions { // XOR
			xor := helpers.Xor(permission.UserEmail, permission.ProgrammaticUserName, permission.TeamId)
			if xor == false {
				resp.Diagnostics.AddError(fmt.Sprintf("Invalid Permission at %s", beautyStringify(permission)), "Exactly one of [user_email, programmatic_username, team_id] must be provided")
			}

			xor = helpers.Xor(permission.Role, permission.CustomRoleId)
			if xor == false {
				resp.Diagnostics.AddError(fmt.Sprintf("Invalid Permission at %s", beautyStringify(permission)), "Exactly one of [role, custom_role_id] must be provided")
			}
		}

		f := func(u *tfNamespacePermissions.PermissionsModel) string {
			return u.GetBlockIdentifier()
		}

		identifiers := helpers.Map(data.Permissions, f)
		if helpers.IsUnique(identifiers) == false {
			duplicates := helpers.FindDuplicates(identifiers, true)
			resp.Diagnostics.AddError("Multiple permissions for the same entity", fmt.Sprintf("'%s' cannot be assigned to multiple permissions", clean(duplicates[0])))
		}

	}
}

// Read refreshes the Terraform state with the latest data.
func (r *NamespacePermissionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get current state
	var state tfNamespacePermissions.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.namespacePermissions.ListNamespacePermissions(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.Diagnostics.AddError("Namespace not found", fmt.Sprintf("Namespace '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read namespace permissions %s", id), fmt.Sprintf("%s", err))
		return
	}

	tfNamespacePermissions.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *NamespacePermissionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan tfNamespacePermissions.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := tfNamespacePermissions.Merge(&plan, nil, commons.CreateMerger)
	namespaceId := plan.NamespaceId

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, namespaceId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = namespaceId

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *NamespacePermissionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan tfNamespacePermissions.ResourceModel
	var state tfNamespacePermissions.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := tfNamespacePermissions.Merge(&plan, &state, commons.UpdateMerger)

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, plan.NamespaceId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//create endpoint is also used for update
	entitiesToUpsert := append(mergeResult.EntitiesToCreate, mergeResult.EntitiesToUpdate...)

	diags = r.createEntities(ctx, entitiesToUpsert, plan.NamespaceId.ValueString())
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

func (r *NamespacePermissionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state tfNamespacePermissions.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := tfNamespacePermissions.Merge(nil, &state, commons.DeleteMerger)

	if diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, state.NamespaceId.ValueString()); diags != nil {
		if diags.ErrorsCount() > 0 && diags[0].Detail() == namespaceNotFoundError {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.Append(diags...)
	}
}

func (r *NamespacePermissionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

//region Private Methods

func (r *NamespacePermissionsResource) createEntities(ctx context.Context, entitiesToCreate []*sdkNamespacePermissions.NamespacePermission, namespaceId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Adding %d permissions to namespace '%s'.", len(entitiesToCreate), namespaceId))

	for _, e := range entitiesToCreate {
		_, err := r.client.Client.namespacePermissions.CreateNamespacePermission(ctx, e)

		if err != nil {
			if commons.DoesErrorContains(err, namespaceNotFoundError) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(
						namespaceNotFoundError,
						fmt.Sprintf("Namespace '%s' Not Found", *e.NamespaceId)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Namespace permission creation failed",
						fmt.Sprintf("failed to create permission for %serror: %s", beautyStringifyApi(e), err)),
				}
			}
		}
	}

	return retVal
}

func (r *NamespacePermissionsResource) deleteEntities(ctx context.Context, entitiesToDelete []*sdkNamespacePermissions.NamespacePermission, namespaceId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Removing %d permissions from namespace '%s'.", len(entitiesToDelete), namespaceId))

	for _, e := range entitiesToDelete {
		partialEntity := &sdkNamespacePermissions.NamespacePermission{
			NamespaceId:          e.NamespaceId,
			UserEmail:            e.UserEmail,
			ProgrammaticUserName: e.ProgrammaticUserName,
			TeamId:               e.TeamId,
		}
		_, err := r.client.Client.namespacePermissions.DeleteNamespacePermission(ctx, partialEntity)

		if err != nil {
			if commons.IsNotFoundResponseError(err) {
				if commons.DoesErrorContains(err, namespaceNotFoundError) {
					return diag.Diagnostics{
						diag.NewErrorDiagnostic("Namespace not found", fmt.Sprintf("Namespace '%s' not found", *e.NamespaceId)),
					}
				} else {
					tflog.Info(ctx, fmt.Sprintf("%sdoes not exists. Removing them from namespace '%s'.", beautyStringifyApi(e), *e.NamespaceId))
				}
			}
		}
	}

	return retVal
}

func clean(s string) string {
	return strings.Split(s, ":")[1]
}

func beautyStringify(e *tfNamespacePermissions.PermissionsModel) string {
	retVal := ""

	if e.UserEmail.IsNull() == false {
		retVal += fmt.Sprintf("user_email: '%s'  ", e.UserEmail.ValueString())
	}
	if e.ProgrammaticUserName.IsNull() == false {
		retVal += fmt.Sprintf("programmatic_username: '%s'  ", e.ProgrammaticUserName.ValueString())
	}
	if e.TeamId.IsNull() == false {
		retVal += fmt.Sprintf("team_id: '%s'  ", e.TeamId.ValueString())
	}
	if e.Role.IsNull() == false {
		retVal += fmt.Sprintf("role: '%s'  ", e.Role.ValueString())
	}
	if e.CustomRoleId.IsNull() == false {
		retVal += fmt.Sprintf("custom_role_id: '%s'  ", e.CustomRoleId.ValueString())
	}

	return retVal
}

func beautyStringifyApi(e *sdkNamespacePermissions.NamespacePermission) string {
	retVal := ""

	if e.UserEmail != nil {
		retVal += fmt.Sprintf("user_email: '%s' ", *e.UserEmail)
	}
	if e.ProgrammaticUserName != nil {
		retVal += fmt.Sprintf("programmatic_username: '%s' ", *e.ProgrammaticUserName)
	}
	if e.TeamId != nil {
		retVal += fmt.Sprintf("team_id: '%s' ", *e.TeamId)
	}
	if e.Role != nil {
		retVal += fmt.Sprintf("role: '%s' ", *e.Role)
	}
	if e.CustomRoleId != nil {
		retVal += fmt.Sprintf("custom_role_id: '%s' ", *e.CustomRoleId)
	}

	return retVal
}

//endregion
