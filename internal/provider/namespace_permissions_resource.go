package provider

import (
	"context"
	"fmt"
	"strings"

	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	sdkNamespacePermissions "github.com/control-monkey/controlmonkey-sdk-go/services/namespace_permissions"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
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
)

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
				MarkdownDescription: "The unique ID of the namespace. Exactly one of  [`namespace_id` , `stack_id`] is required.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("stack_id"), path.MatchRoot("namespace_id")),
				},
			},
			"stack_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the stack. Exactly one of  [`namespace_id` , `stack_id`] is required.",
				Optional:            true,
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

	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		return
	}

	namespaceIdStackIdXor := helpers.Xor(data.NamespaceId, data.StackId)

	if namespaceIdStackIdXor == false {
		resp.Diagnostics.AddError(validationError, "Exactly one of [namespace_id, stack_id] must be provided.")
		return
	}

	if len(data.Permissions) > 0 {
		//Validate constraints
		for _, permission := range data.Permissions { // XOR
			xor := helpers.Xor(permission.UserEmail, permission.ProgrammaticUserName, permission.TeamId)
			if xor == false {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("Invalid Permission at %s. Exactly one of [user_email, programmatic_username, team_id] must be provided", stringifyError(permission)))
				return
			}

			xor = helpers.Xor(permission.Role, permission.CustomRoleId)
			if xor == false {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("Invalid Permission at %s. Exactly one of [role, custom_role_id] must be provided.", stringifyError(permission)))
				return
			}
		}

		identifiers := interfaces.GetIdentifiers(data.Permissions)

		if helpers.IsUnique(identifiers) == false {
			duplicates := helpers.FindDuplicates(identifiers, false)
			for _, d := range duplicates {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("'%s' cannot be assigned to multiple permissions", tfNamespacePermissions.CleanIdentifier(d)))
			}
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

	var namespaceId *string
	var stackId *string

	if state.NamespaceId.IsNull() == false {
		namespaceId = &id
	} else if state.StackId.IsNull() == false {
		stackId = &id
	}

	res, err := r.client.Client.namespacePermissions.ListNamespacePermissions(ctx, namespaceId, stackId)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(namespaceNotFoundError, fmt.Sprintf("Namespace or Stack '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read namespace permissions for namespace '%s'", id), err.Error())
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

	// ID can be either namespaceId or stackId
	var id string
	if plan.NamespaceId.IsNull() == false {
		id = plan.NamespaceId.ValueString()
	} else {
		id = plan.StackId.ValueString()
	}

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(id)

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

	// ID can be either namespaceId or stackId
	var id string
	if plan.NamespaceId.IsNull() == false {
		id = plan.NamespaceId.ValueString()
	} else {
		id = plan.StackId.ValueString()
	}

	mergeResult := tfNamespacePermissions.Merge(&plan, &state, commons.UpdateMerger)

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//create endpoint is also used for update
	entitiesToUpsert := append(mergeResult.EntitiesToCreate, mergeResult.EntitiesToUpdate...)

	diags = r.createEntities(ctx, entitiesToUpsert, id)
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

	// ID can be either namespaceId or stackId
	var id string
	if state.NamespaceId.IsNull() == false {
		id = state.NamespaceId.ValueString()
	} else {
		id = state.StackId.ValueString()
	}

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, id)
	resp.Diagnostics.Append(diags...)
}

func (r *NamespacePermissionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID

	// ID can be either namespaceId or stackId, Identify which one to use based on the prefix
	if strings.HasPrefix(id, "ns-") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("namespace_id"), id)...)
	} else if strings.HasPrefix(id, "stk-") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("stack_id"), id)...)
	} else {
		// Unknown format
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be either a namespace ID (ns-xxx) or stack ID (stk-xxx), got: %s", id),
		)
		return
	}
}

//region Private Methods

func (r *NamespacePermissionsResource) createEntities(ctx context.Context, entitiesToCreate []*sdkNamespacePermissions.NamespacePermission, id string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Adding %d permissions to namespace/stack '%s'.", len(entitiesToCreate), id))

	for _, e := range entitiesToCreate {
		_, err := r.client.Client.namespacePermissions.CreateNamespacePermission(ctx, e)

		if err != nil {
			if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to add permission '%s' to namespace/stack '%s'. Error: %s",
						beautyStringifyApi(e), id, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to add permission '%s' to namespace/stack '%s'", beautyStringifyApi(e), id),
						err.Error()),
				}
			}
		}
	}

	return retVal
}

func (r *NamespacePermissionsResource) deleteEntities(ctx context.Context, entitiesToDelete []*sdkNamespacePermissions.NamespacePermission, id string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Removing %d permissions from namespace/stack '%s'.", len(entitiesToDelete), id))

	for _, e := range entitiesToDelete {
		partialEntity := &sdkNamespacePermissions.NamespacePermission{
			NamespaceId:          e.NamespaceId,
			StackId:              e.StackId,
			UserEmail:            e.UserEmail,
			ProgrammaticUserName: e.ProgrammaticUserName,
			TeamId:               e.TeamId,
		}
		_, err := r.client.Client.namespacePermissions.DeleteNamespacePermission(ctx, partialEntity)

		if err != nil {
			if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to remove permission '%s' from namespace/stack '%s'. Error: %s",
						beautyStringifyApi(e), id, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to remove permission '%s' from namespace/stack '%s'", beautyStringifyApi(e), id),
						err.Error()),
				}
			}
		}
	}

	return retVal
}

func stringifyError(e *tfNamespacePermissions.PermissionsModel) string {
	retVal := ""
	retVal += stringifyProperty(e.UserEmail, "user_email")
	retVal += stringifyProperty(e.ProgrammaticUserName, "programmatic_username")
	retVal += stringifyProperty(e.TeamId, "team_id")
	retVal += stringifyProperty(e.Role, "role")
	retVal += stringifyProperty(e.CustomRoleId, "custom_role_id")

	return retVal
}

func stringifyProperty(e types.String, propertyName string) string {
	retVal := ""

	if e.IsNull() == false {
		if e.IsUnknown() {
			retVal = "'known after apply'"
		} else {
			retVal += fmt.Sprintf("%s: '%s'  ", propertyName, e.ValueString())
		}
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
