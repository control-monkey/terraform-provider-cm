package provider

import (
	"context"
	"fmt"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	tfCustomRole "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/custom_role"
	cm_stringvalidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
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
var _ resource.Resource = &CustomRoleResource{}

func NewCustomRoleResource() resource.Resource {
	return &CustomRoleResource{}
}

type CustomRoleResource struct {
	client *ControlMonkeyAPIClient
}

func (r *CustomRoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_role"
}

func (r *CustomRoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys custom roles.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the custom role.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				// Not Computed: the API defaults to namespaceRole but never returns it, so a Computed
				// value would stay unknown after apply.
				MarkdownDescription: fmt.Sprintf("The type of the role. Allowed values: %s."+
					" **Omitting this creates a `%s`.** This attribute cannot be modified after creation.",
					helpers.EnumForDocs(cmTypes.CustomRoleTypeTypes), cmTypes.NamespaceRoleType),
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(cmTypes.CustomRoleTypeTypes...),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the role.",
				Required:            true,
				Validators: []validator.String{
					cm_stringvalidators.NotBlank(),
					stringvalidator.NoneOf("admin", "viewer"),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the role.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.NoneOf(""),
				},
			},
			"permissions": schema.ListNestedAttribute{
				MarkdownDescription: "List of permissions allowed by the role.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: fmt.Sprintf("The type of the permission. Only allowed when `type` is `%s`."+
								" Find supported types [here](https://docs.controlmonkey.io/controlmonkey-api/api-enumerations#custom-role-permission-types).", cmTypes.NamespaceRoleType),
							Optional: true,
						},
						"names": schema.ListAttribute{
							MarkdownDescription: fmt.Sprintf("The types of the permissions. Only allowed when `type` is `%s`."+
								" Find supported types [here](https://docs.controlmonkey.io/controlmonkey-api/api-enumerations#custom-role-permission-types).", cmTypes.OrganizationRoleType),
							ElementType: types.StringType,
							Optional:    true,
						},
						"restrictions": schema.ListNestedAttribute{
							MarkdownDescription: fmt.Sprintf("Narrows the permissions to matching resources only. Only allowed when `type` is `%s`,"+
								" and only for permissions that support restrictions."+
								" Fields within one restriction are combined with AND; multiple restrictions are combined with OR."+
								" All fields accept regular expressions, matched case-sensitively.", cmTypes.OrganizationRoleType),
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"cloud_provider": schema.StringAttribute{
										MarkdownDescription: "Restrict to matching cloud providers, for example `aws`. Must be lowercase.",
										Optional:            true,
									},
									"cloud_account_id": schema.StringAttribute{
										MarkdownDescription: "Restrict to matching cloud account IDs.",
										Optional:            true,
									},
									"cm_resource_name": schema.StringAttribute{
										MarkdownDescription: "Restrict to matching ControlMonkey resource names.",
										Optional:            true,
									},
								},
							},
						},
					},
				},
			},
			"stack_restriction": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("Restrict stack operations with supported types. Only allowed when `type` is `%s`."+
					" Learn more [here](https://docs.controlmonkey.io/administration/users-and-roles/custom-roles)."+
					" Find supported types [here](https://docs.controlmonkey.io/controlmonkey-api/api-enumerations#stack-restriction-types).", cmTypes.NamespaceRoleType),
				Optional: true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *CustomRoleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ValidateConfig enforces the attributes that are only valid for one role type. Action names are
// left to the API.
func (r *CustomRoleResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data tfCustomRole.ResourceModel

	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		return
	}

	// An unset type defaults to namespaceRole on the API side. An unknown type matches neither.
	isOrganizationRole := data.Type.ValueString() == cmTypes.OrganizationRoleType
	isNamespaceRole := isOrganizationRole == false && data.Type.IsUnknown() == false

	if isOrganizationRole && data.StackRestriction.IsNull() == false {
		resp.Diagnostics.AddError(validationError, fmt.Sprintf("stack_restriction is allowed only when type is '%s'", cmTypes.NamespaceRoleType))
	}

	for _, p := range data.Permissions {
		if isOrganizationRole {
			if p.Name.IsNull() == false {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("permissions.name is not allowed when type is '%s'. Use permissions.names instead", cmTypes.OrganizationRoleType))
			}
		} else if isNamespaceRole {
			if p.Names != nil {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("permissions.names is allowed only when type is '%s'. Use permissions.name instead", cmTypes.OrganizationRoleType))
			}
			if p.Restrictions != nil {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("permissions.restrictions is allowed only when type is '%s'", cmTypes.OrganizationRoleType))
			}
		}

		if p.Name.IsNull() && p.Names == nil {
			resp.Diagnostics.AddError(validationError, "Exactly one of [permissions.name, permissions.names] must be set")
		}
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *CustomRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state tfCustomRole.ResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.customRole.ReadCustomRole(ctx, id)

	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(resourceNotFoundError, fmt.Sprintf("Custom role '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read custom role '%s'", id), err.Error())
		return
	}

	tfCustomRole.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *CustomRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan tfCustomRole.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := tfCustomRole.Converter(&plan, nil, commons.CreateConverter)

	res, err := r.client.Client.customRole.CreateCustomRole(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceCreationFailedError,
			fmt.Sprintf("failed to create custom role, error: %s", err.Error()),
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

func (r *CustomRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan tfCustomRole.ResourceModel
	var state tfCustomRole.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body, _ := tfCustomRole.Converter(&plan, &state, commons.UpdateConverter)

	_, err := r.client.Client.customRole.UpdateCustomRole(ctx, id, body)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.Diagnostics.AddError(resourceNotFoundError, fmt.Sprintf("Custom role '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(
			resourceUpdateFailedError,
			fmt.Sprintf("failed to update custom role %s, error: %s", id, err),
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

func (r *CustomRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state tfCustomRole.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	_, err := r.client.Client.customRole.DeleteCustomRole(ctx, id)

	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			resourceDeletionFailedError,
			fmt.Sprintf("Failed to delete custom role %s, error: %s", id, err),
		)
		return
	}
}

func (r *CustomRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
