package provider

import (
	"context"
	"fmt"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	sdkControlPolicy "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	controlPolicyMapping "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/control_policy_mapping"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ControlPolicyMappingResource{}

func NewControlPolicyMappingResource() resource.Resource {
	return &ControlPolicyMappingResource{}
}

type ControlPolicyMappingResource struct {
	client *ControlMonkeyAPIClient
}

func (r *ControlPolicyMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_control_policy_mapping"
}

func (r *ControlPolicyMappingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys Control Policy Mappings.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"control_policy_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the Control Policy.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"target_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID corresponds to the `target_type` in the mapping.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"target_type": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The type of the target. Allowed values: %s.", helpers.EnumForDocs(cmTypes.PolicyMappingTargetTypes)),
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(cmTypes.PolicyMappingTargetTypes...),
				},
			},
			"enforcement_level": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("Specifies the level of enforcement for the control policy on the target. Allowed values: %s."+
					" When set to `softMandatory`, a policy failure triggers an approval requirement before applying changes."+
					" When set to `hardMandatory`, changes cannot be applied until the policy check is successful.", helpers.EnumForDocs(cmTypes.EnforcementLevelTypes)),
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(cmTypes.EnforcementLevelTypes...),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *ControlPolicyMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ControlPolicyMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state controlPolicyMapping.ResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	controlPolicyId := state.ControlPolicyId.ValueString()
	targetId := state.TargetId.ValueString()
	targetType := state.TargetType.ValueString()

	res, err := r.client.Client.controlPolicyMapping.ListControlPolicyMappings(ctx, controlPolicyId)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read controlPolicyMapping for control_policy_id '%s' and target_id '%s' of type '%s'", controlPolicyId, targetId, targetType), fmt.Sprintf("%s", err))
		return
	}

	filter := func(m *sdkControlPolicy.ControlPolicyMapping) bool {
		cond1 := *m.ControlPolicyId == state.ControlPolicyId.ValueString()
		cond2 := *m.TargetId == state.TargetId.ValueString()
		cond3 := *m.TargetType == state.TargetType.ValueString()

		return cond1 && cond2 && cond3
	}

	mapping := helpers.FindFirst(res, filter)

	if mapping == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	controlPolicyMapping.UpdateStateAfterRead(mapping, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ControlPolicyMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan controlPolicyMapping.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := controlPolicyMapping.Converter(&plan, nil, commons.CreateConverter)

	_, err := r.client.Client.controlPolicyMapping.CreateControlPolicyMapping(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"ControlPolicyMapping creation failed",
			fmt.Sprintf("failed to create controlPolicyMapping, error: %s", err.Error()),
		)
		return
	}

	plan.ID = controlPolicyMapping.ComputeId(body)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ControlPolicyMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan controlPolicyMapping.ResourceModel
	var state controlPolicyMapping.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := controlPolicyMapping.Converter(&plan, &state, commons.UpdateConverter)

	_, err := r.client.Client.controlPolicyMapping.UpdateControlPolicyMapping(ctx, body)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"ControlPolicyMapping update failed",
			fmt.Sprintf("failed to update controlPolicyMapping, error: %s", err.Error()),
		)
		return
	}

	plan.ID = controlPolicyMapping.ComputeId(body)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ControlPolicyMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state controlPolicyMapping.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Client.controlPolicyMapping.DeleteControlPolicyMapping(ctx, &sdkControlPolicy.ControlPolicyMapping{
		ControlPolicyId: state.ControlPolicyId.ValueStringPointer(),
		TargetId:        state.TargetId.ValueStringPointer(),
		TargetType:      state.TargetType.ValueStringPointer(),
	})

	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"ControlPolicyMapping deletion failed",
			fmt.Sprintf("Failed to delete controlPolicyMapping, error: %s", err),
		)
		return
	}
}

func (r *ControlPolicyMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")
	anyBlank := helpers.AnyMatch(idParts, helpers.IsBlank)

	if len(idParts) != 3 || anyBlank {
		resp.Diagnostics.AddError("Import cm_control_policy_mapping Failed",
			fmt.Sprintf("Unexpected format of ID (%q), expected control_policy_id/target_id/target_type", req.ID))
		return
	}

	state := new(controlPolicyMapping.ResourceModel)
	state.ControlPolicyId = helpers.StringValueOrNull(&idParts[0])
	state.TargetId = helpers.StringValueOrNull(&idParts[1])
	state.TargetType = helpers.StringValueOrNull(&idParts[2])

	resp.State.Set(ctx, &state)
}
