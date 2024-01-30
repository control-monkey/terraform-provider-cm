package provider

import (
	"context"
	"fmt"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	sdkVariable "github.com/control-monkey/controlmonkey-sdk-go/services/variable"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/variable"
	cm_listvalidator "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/list"
	cm_stringvalidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
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

func (r *VariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

func (r *VariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys variables.\nVariable can be either a Terraform variable or an Environment variable.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the variable.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("Scope of the variable. Allowed values: %s.", helpers.EnumForDocs(cmTypes.VariableScopeTypes)),
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(cmTypes.VariableScopeTypes...),
				},
			},
			"scope_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the resource to which the variable is attached.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "The key of the variable.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The value of the variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("Type of the variable. Allowed values: %s.", helpers.EnumForDocs(cmTypes.VariableTypes)),
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(cmTypes.VariableTypes...),
				},
			},
			"is_sensitive": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the variable value is sensitive and requires encryption.",
				Required:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"is_overridable": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the variable can be overridden by a lower-level scope.",
				Required:            true,
			},
			"is_required": schema.BoolAttribute{
				MarkdownDescription: "This setting applies to template variables without a specified value. Stacks created from the template need to provide a value for this variable.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for the variable.",
				Optional:            true,
			},
			"value_conditions": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Specify conditions for the variable value using an operator and another value. Typically used for stacks launched from templates. For more information: [ControlMonkey Docs] (https://docs.controlmonkey.io/main-concepts/variables/variable-conditions)",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"operator": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: fmt.Sprintf("Logical operators. Allowed values: %s.", helpers.EnumForDocs(cmTypes.VariableConditionOperatorTypes)),
							Validators: []validator.String{
								stringvalidator.OneOf(cmTypes.VariableConditionOperatorTypes...),
							},
						},
						"value": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: fmt.Sprintf("The value associated with the operator. Input a number or string depending on the chosen operator. Use `values` field for operator of type `%s`", cmTypes.In),
							Validators:          []validator.String{cm_stringvalidators.NotBlank(), stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("values"))},
						},
						"values": schema.ListAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							MarkdownDescription: fmt.Sprintf("A list of strings when using operator type `%s`. For other operators use `value`", cmTypes.In),
							Validators:          commons.ValidateUniqueNotEmptyListWithNoBlankValues(),
						},
					},
				},
				Validators: []validator.List{
					cm_listvalidator.SizeExactly(1),
				},
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

func (r *VariableResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data variable.ResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	valueConditions := data.ValueConditions

	if valueConditions != nil {
		for i, condition := range valueConditions {
			var errMsg string

			switch op := condition.Operator.ValueString(); op {
			case cmTypes.Ne:
				if condition.Value.ValueStringPointer() == nil {
					errMsg = fmt.Sprintf("value_conditions[%d].value must be set", i)
				}
			case cmTypes.Gt, cmTypes.Gte, cmTypes.Lt, cmTypes.Lte:
				isNumeric, _ := helpers.CheckAndGetIfNumericString(condition.Value.ValueString())
				if !isNumeric {
					errMsg = fmt.Sprintf("value_conditions[%d].value must be a number when using value_conditions.operator '%s'", i, op)
				}
			case cmTypes.In:
				if condition.Values.IsNull() {
					errMsg = fmt.Sprintf("value_conditions[%d].values must be set", i)
				}
			case cmTypes.StartsWith, cmTypes.Contains:
				if condition.Value.ValueStringPointer() == nil {
					errMsg = fmt.Sprintf("value_conditions[%d].value must be set", i)
				}
			}

			if errMsg != "" {
				resp.Diagnostics.AddError("Validation Error", errMsg)
			}
		}
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *VariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state variable.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.variable.ReadVariable(ctx, &sdkVariable.ReadVariableInput{VariableId: controlmonkey.String(id)})
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read variable %s", id), fmt.Sprintf("%s", err))
		return
	}

	variable.UpdateStateAfterRead(res, &state)

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
	var plan variable.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := variable.Converter(&plan, nil, commons.CreateConverter)

	res, err := r.client.Client.variable.CreateVariable(ctx, body)
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
	var plan variable.ResourceModel
	var state variable.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueStringPointer()

	body, _ := variable.Converter(&plan, &state, commons.UpdateConverter)

	_, err := r.client.Client.variable.UpdateVariable(ctx, id, body)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Variable update failed",
			fmt.Sprintf("failed to update variable %s, error: %s", *id, err),
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

func (r *VariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state variable.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueStringPointer()

	_, err := r.client.Client.variable.DeleteVariable(ctx, &sdkVariable.DeleteVariableInput{
		VariableId: id,
	})
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Variable deletion failed",
			fmt.Sprintf("Failed to delete variable %s, error: %s", *id, err),
		)
		return
	}
}

func (r *VariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
