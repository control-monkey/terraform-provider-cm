package provider

import (
	"context"
	"fmt"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	sdkControlPolicyGroup "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy_group"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	controlPolicyGroupMapping "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/control_policy_group_mappings"
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
var _ resource.Resource = &ControlPolicyGroupMappingResource{}

func NewControlPolicyGroupMappingResource() resource.Resource {
	return &ControlPolicyGroupMappingResource{}
}

type ControlPolicyGroupMappingResource struct {
	client *ControlMonkeyAPIClient
}

func (r *ControlPolicyGroupMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_control_policy_group_mappings"
}

func (r *ControlPolicyGroupMappingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys control policy group mappings.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"control_policy_group_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the control policy group.",
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
							MarkdownDescription: fmt.Sprintf("Specifies the level of enforcement for the control policy group on the target. Allowed values: %s."+
								" When set to `softMandatory`, a policy failure triggers an approval requirement before applying changes."+
								" When set to `hardMandatory`, changes cannot be applied until the policy check is successful."+
								" When set to `bySeverity`, the enforcement level will be the determined by the severity of each policy in the policy group", helpers.EnumForDocs(cmTypes.GroupEnforcementLevelTypes)),
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(cmTypes.GroupEnforcementLevelTypes...),
							},
						},
						"override_enforcements": schema.SetNestedAttribute{
							MarkdownDescription: fmt.Sprintf(""),
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"control_policy_id": schema.StringAttribute{
										MarkdownDescription: "The unique ID of the control policy.",
										Required:            true,
										Validators: []validator.String{
											cmStringValidators.NotBlank(),
										},
									},
									"enforcement_level": schema.StringAttribute{
										MarkdownDescription: fmt.Sprintf("Specifies the level of enforcement for the control policy group on the target. Allowed values: %s."+
											" When set to `softMandatory`, a policy failure triggers an approval requirement before applying changes."+
											" When set to `hardMandatory`, changes cannot be applied until the policy check is successful.", helpers.EnumForDocs(cmTypes.EnforcementLevelTypes)),
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf(cmTypes.EnforcementLevelTypes...),
										},
									},
									"stack_ids": schema.ListAttribute{
										MarkdownDescription: fmt.Sprintf("A list of stack IDs within the specified namespace where the original enforcement level of the policy will be overridden with the new enforcement level. This option can only be used when the `target_type` is set to '%s'.", cmTypes.NamespaceTargetType),
										ElementType:         types.StringType,
										Optional:            true,
										Validators:          commons.ValidateUniqueNotEmptyListWithNoBlankValues(),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *ControlPolicyGroupMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ControlPolicyGroupMappingResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data controlPolicyGroupMapping.ResourceModel

	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		return
	}

	if len(data.Targets) > 0 {
		identifiers := interfaces.GetIdentifiers(data.Targets)
		if helpers.IsUnique(identifiers) == false {
			duplicates := helpers.FindDuplicates(identifiers, false)
			for _, d := range duplicates {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("Target '%s' appears more than once", controlPolicyGroupMapping.CleanTargetIdentifier(d)))
			}

			if resp.Diagnostics.HasError() {
				return
			}
		}

		for _, t := range data.Targets {
			if targetType := t.TargetType; helpers.IsKnown(targetType) {
				if len(t.OverrideEnforcements) > 0 {
					for _, o := range t.OverrideEnforcements {
						if helpers.IsKnown(o.StackIds) && targetType.ValueString() == cmTypes.StackTargetType {
							resp.Diagnostics.AddError(validationError, fmt.Sprintf("Target type '%s' cannot have stack_ids set", targetType.ValueString()))
							return
						}
					}
				}
			}
		}
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *ControlPolicyGroupMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state controlPolicyGroupMapping.ResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.controlPolicyGroup.ListControlPolicyGroupMappings(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(controlPolicyGroupNotFoundError, fmt.Sprintf("Control policy group '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read mappings for control policy group '%s'", id), err.Error())
		return
	}

	controlPolicyGroupMapping.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ControlPolicyGroupMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan controlPolicyGroupMapping.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := controlPolicyGroupMapping.Merge(&plan, nil, commons.CreateConverter)
	controlPolicyGroupId := plan.ControlPolicyGroupId

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, controlPolicyGroupId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = controlPolicyGroupId

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ControlPolicyGroupMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan controlPolicyGroupMapping.ResourceModel
	var state controlPolicyGroupMapping.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := controlPolicyGroupMapping.Merge(&plan, &state, commons.UpdateMerger)
	controlPolicyGroupId := plan.ControlPolicyGroupId

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, controlPolicyGroupId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.updateEntities(ctx, mergeResult.EntitiesToUpdate, controlPolicyGroupId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, controlPolicyGroupId.ValueString())
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

func (r *ControlPolicyGroupMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state controlPolicyGroupMapping.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := controlPolicyGroupMapping.Merge(nil, &state, commons.DeleteMerger)
	controlPolicyGroupId := state.ControlPolicyGroupId

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, controlPolicyGroupId.ValueString())
	resp.Diagnostics.Append(diags...)
}

func (r *ControlPolicyGroupMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ControlPolicyGroupMappingResource) createEntities(ctx context.Context, entitiesToCreate []*sdkControlPolicyGroup.ControlPolicyGroupMapping, controlPolicyGroupId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Mapping %d targets to Control Policy Group '%s'.", len(entitiesToCreate), controlPolicyGroupId))

	for _, e := range entitiesToCreate {
		_, err := r.client.Client.controlPolicyGroup.CreateControlPolicyGroupMapping(ctx, e)

		if err != nil {
			targetId := *e.TargetId
			targetType := *e.TargetType
			if commons.IsAlreadyExistResponseError(err) {
				tflog.Info(ctx, fmt.Sprintf("Target '%s' of type '%s' is already mapped to control policy group '%s'. No operation was made.", targetId, targetType, controlPolicyGroupId))
			} else if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to create map between target '%s' of type '%s' and control policy group '%s'. Error: %s", targetType, targetId, controlPolicyGroupId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to create map between target '%s' of type '%s' and control policy group '%s'", targetType, targetId, controlPolicyGroupId),
						err.Error()),
				}
			}
		}
	}

	return retVal
}

func (r *ControlPolicyGroupMappingResource) updateEntities(ctx context.Context, entitiesToUpdate []*sdkControlPolicyGroup.ControlPolicyGroupMapping, controlPolicyGroupId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Updating %d target mappings to Control Policy Group '%s'.", len(entitiesToUpdate), controlPolicyGroupId))

	for _, e := range entitiesToUpdate {
		_, err := r.client.Client.controlPolicyGroup.UpdateControlPolicyGroupMapping(ctx, e)

		if err != nil {
			targetId := *e.TargetId
			targetType := *e.TargetType
			if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to update map between target '%s' of type '%s' and control policy group '%s'. Error: %s", targetType, targetId, controlPolicyGroupId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to update map between target '%s' of type '%s' and control policy group '%s'", targetType, targetId, controlPolicyGroupId),
						err.Error()),
				}
			}
		}
	}

	return retVal
}

func (r *ControlPolicyGroupMappingResource) deleteEntities(ctx context.Context, entitiesToDelete []*sdkControlPolicyGroup.ControlPolicyGroupMapping, controlPolicyGroupId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Removing %d target mappings from control policy group '%s'.", len(entitiesToDelete), controlPolicyGroupId))

	for _, e := range entitiesToDelete {
		_, err := r.client.Client.controlPolicyGroup.DeleteControlPolicyGroupMapping(ctx, e)

		if err != nil {
			targetId := *e.TargetId
			targetType := *e.TargetType
			if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to delete map between target '%s' of type '%s' and control policy group '%s'. Error: %s", targetType, targetId, controlPolicyGroupId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to delete map between target '%s' of type '%s' and control policy group '%s'", targetType, targetId, controlPolicyGroupId),
						err.Error()),
				}
			}
		}
	}

	return retVal
}
