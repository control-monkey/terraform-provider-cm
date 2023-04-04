package provider

import (
	"context"
	"fmt"
	"github.com/control-monkey-customer-z/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	"github.com/control-monkey/controlmonkey-sdk-go/service/variable"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &VariableResource{}

func NewVariableResource() resource.Resource {
	return &VariableResource{}
}

type VariableResource struct {
	client *ControlMonkeyAPIClient
}

type VariableResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Scope         types.String `tfsdk:"scope"`
	ScopeId       types.String `tfsdk:"scope_id"`
	Key           types.String `tfsdk:"key"`
	Type          types.String `tfsdk:"type"`
	Value         types.String `tfsdk:"value"`
	IsSensitive   types.Bool   `tfsdk:"is_sensitive"`
	IsOverridable types.Bool   `tfsdk:"is_overridable"`
	IsRequired    types.Bool   `tfsdk:"is_required"`
	Description   types.String `tfsdk:"description"`
}

var vScopes = []string{
	"organization",
	"namespace",
	"variableGroup",
	"template",
	"stack",
	"stackRun",
}

var vTypes = []string{
	"tfVar",
	"envVar",
}

func (r *VariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

func (r *VariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Variable can be either a Terraform variable or an Environment variable.\n Terraform variables, which are defined in the Terraform configuration files and have their values set by ControlMonkey when running IaC commands. \nEnvironment variables, which are set by ControlMonkey in the shell running IaC commands",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the variable",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("Scope of the variable. Allowed values: %v", vScopes),
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(vScopes...),
				},
			},
			"scope_id": schema.StringAttribute{
				MarkdownDescription: "The id of the scope resource that the variable is attached to",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "The key of the variable",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The value of the variable",
				Optional:            true,
				Sensitive:           true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("Type of the variable. Allowed values: %v", vScopes),
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(vTypes...),
				},
			},
			"is_sensitive": schema.BoolAttribute{
				MarkdownDescription: "Whether the variable value is sensitive and should be encrypted or not",
				Required:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"is_overridable": schema.BoolAttribute{
				MarkdownDescription: "Either the value of the variable can be overridden by a another scope down in the hierarchy or the variable will be inherited without modifications",
				Required:            true,
			},
			"is_required": schema.BoolAttribute{
				MarkdownDescription: "This setting is used for template variables that do not have a value specified. Stacks created from the template must assign a value to this variable.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for the variable",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *VariableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *VariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state VariableResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.variable.ReadVariable(ctx, &variable.ReadVariableInput{VariableId: controlmonkey.String(id)})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read variable %s", id), fmt.Sprintf("%s", err))
		return
	}

	respVariable := res.Variable

	state.Scope = helpers.StringValueOrNull(respVariable.Scope)
	state.Key = helpers.StringValueOrNull(respVariable.Key)
	state.Type = helpers.StringValueOrNull(respVariable.Type)
	state.IsSensitive = helpers.BoolValueOrNull(respVariable.IsSensitive)
	state.IsOverridable = helpers.BoolValueOrNull(respVariable.IsOverridable)
	state.ScopeId = helpers.StringValueOrNull(respVariable.ScopeId)
	state.Value = helpers.StringValueOrNull(respVariable.Value)
	state.IsRequired = helpers.BoolValueOrNull(respVariable.IsRequired)
	state.Description = helpers.StringValueOrNull(respVariable.Description)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *VariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan VariableResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body variable.Variable

	body.SetScope(plan.Scope.ValueStringPointer())
	body.SetScopeId(plan.ScopeId.ValueStringPointer())
	body.SetKey(plan.Key.ValueStringPointer())
	body.SetValue(plan.Value.ValueStringPointer())
	body.SetType(plan.Type.ValueStringPointer())
	body.SetIsSensitive(plan.IsSensitive.ValueBoolPointer())
	body.SetIsOverridable(plan.IsOverridable.ValueBoolPointer())
	body.SetIsRequired(plan.IsRequired.ValueBoolPointer())
	body.SetDescription(plan.Description.ValueStringPointer())

	res, err := r.client.Client.variable.CreateVariable(ctx, &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Variable creation failed",
			fmt.Sprintf("failed to create variable, error: %s", err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(controlmonkey.StringValue(res.Variable.ID))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *VariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan VariableResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueStringPointer()

	var body variable.Variable

	body.SetValue(plan.Value.ValueStringPointer())
	body.SetIsOverridable(plan.IsOverridable.ValueBoolPointer())
	body.SetIsRequired(plan.IsRequired.ValueBoolPointer())
	body.SetDescription(plan.Description.ValueStringPointer())

	_, err := r.client.Client.variable.UpdateVariable(ctx, id, &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Variable update failed",
			fmt.Sprintf("failed to update variable %s, error: %s", *id, err.Error()),
		)
		return
	}

	// NOTE: no need to update ID and ProjectID
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *VariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state VariableResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueStringPointer()

	_, err := r.client.Client.variable.DeleteVariable(ctx, &variable.DeleteVariableInput{
		VariableId: id,
	})

	if err != nil {
		errMsg := err.Error()
		resp.Diagnostics.AddError(
			"Variable deletion failed",
			fmt.Sprintf("Failed to delete variable %s, error: %s", *id, errMsg),
		)
		return
	}
}

func (r *VariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
