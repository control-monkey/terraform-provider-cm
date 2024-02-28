package provider

import (
	"context"
	"fmt"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	sdkControlPolicy "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	controlPolicyMapping "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/control_policy_mappings"
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
var _ resource.Resource = &ControlPolicyMappingResource{}

func NewControlPolicyMappingResource() resource.Resource {
	return &ControlPolicyMappingResource{}
}

type ControlPolicyMappingResource struct {
	client *ControlMonkeyAPIClient
}

func (r *ControlPolicyMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_control_policy_mappings"
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
			"targets": schema.SetNestedAttribute{
				MarkdownDescription: "List of targets",
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
						"target_id": schema.StringAttribute{
							MarkdownDescription: "The unique ID corresponds to the `target_type` in the mapping.",
							Required:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
						"target_type": schema.StringAttribute{
							MarkdownDescription: fmt.Sprintf("The type of the target. Allowed values: %s.", helpers.EnumForDocs(cmTypes.PolicyMappingTargetTypes)),
							Required:            true,
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

func (r *ControlPolicyMappingResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data controlPolicyMapping.ResourceModel

	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		return
	}

	if len(data.Targets) > 0 {
		identifiers := interfaces.GetIdentifiers(data.Targets)

		if helpers.IsUnique(identifiers) == false {
			duplicates := helpers.FindDuplicates(identifiers, false)
			for _, d := range duplicates {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("Target '%s' appears more than once", controlPolicyMapping.CleanIdentifier(d)))
			}
		}
	}
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

	id := state.ID.ValueString()
	res, err := r.client.Client.controlPolicy.ListControlPolicyMappings(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(controlPolicyNotFoundError, fmt.Sprintf("Control policy '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read mappings for control policy '%s'", id), err.Error())
		return
	}

	controlPolicyMapping.UpdateStateAfterRead(res, &state)

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

	mergeResult := controlPolicyMapping.Merge(&plan, nil, commons.CreateConverter)
	controlPolicyId := plan.ControlPolicyId

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, controlPolicyId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = controlPolicyId

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

	mergeResult := controlPolicyMapping.Merge(&plan, &state, commons.UpdateMerger)
	controlPolicyId := plan.ControlPolicyId

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, controlPolicyId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.updateEntities(ctx, mergeResult.EntitiesToUpdate, controlPolicyId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, controlPolicyId.ValueString())
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

func (r *ControlPolicyMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state controlPolicyMapping.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := controlPolicyMapping.Merge(nil, &state, commons.DeleteMerger)
	controlPolicyId := state.ControlPolicyId

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, controlPolicyId.ValueString())
	resp.Diagnostics.Append(diags...)
}

func (r *ControlPolicyMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ControlPolicyMappingResource) createEntities(ctx context.Context, entitiesToCreate []*sdkControlPolicy.ControlPolicyMapping, controlPolicyId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Mapping %d targets to control policy '%s'.", len(entitiesToCreate), controlPolicyId))

	for _, e := range entitiesToCreate {
		_, err := r.client.Client.controlPolicy.CreateControlPolicyMapping(ctx, e)

		if err != nil {
			targetId := *e.TargetId
			targetType := *e.TargetType
			if commons.IsAlreadyExistResponseError(err) {
				tflog.Info(ctx, fmt.Sprintf("Target '%s' of type '%s' is already mapped to control policy '%s'. No operation was made.", targetId, targetType, controlPolicyId))
			} else if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to create map between target '%s' of type '%s' and control policy '%s'. Error: %s", targetType, targetId, controlPolicyId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to create map between target '%s' of type '%s' and control policy '%s'", targetType, targetId, controlPolicyId),
						err.Error()),
				}
			}
		}
	}

	return retVal
}

func (r *ControlPolicyMappingResource) updateEntities(ctx context.Context, entitiesToUpdate []*sdkControlPolicy.ControlPolicyMapping, controlPolicyId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Updating %d target mappings to Control Policy '%s'.", len(entitiesToUpdate), controlPolicyId))

	for _, e := range entitiesToUpdate {
		_, err := r.client.Client.controlPolicy.UpdateControlPolicyMapping(ctx, e)

		if err != nil {
			targetId := *e.TargetId
			targetType := *e.TargetType
			if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to update map between target '%s' of type '%s' and control policy '%s'. Error: %s", targetType, targetId, controlPolicyId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to update map between target '%s' of type '%s' and control policy '%s'", targetType, targetId, controlPolicyId),
						err.Error()),
				}
			}
		}
	}

	return retVal
}

func (r *ControlPolicyMappingResource) deleteEntities(ctx context.Context, entitiesToDelete []*sdkControlPolicy.ControlPolicyMapping, controlPolicyId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Removing %d target mappings from control policy '%s'.", len(entitiesToDelete), controlPolicyId))

	for _, e := range entitiesToDelete {
		_, err := r.client.Client.controlPolicy.DeleteControlPolicyMapping(ctx, e)

		if err != nil {
			targetId := *e.TargetId
			targetType := *e.TargetType
			if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to delete map between target '%s' of type '%s' and control policy '%s'. Error: %s", targetType, targetId, controlPolicyId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to delete map between target '%s' of type '%s' and control policy '%s'", targetType, targetId, controlPolicyId),
						err.Error()),
				}
			}
		}
	}

	return retVal
}
