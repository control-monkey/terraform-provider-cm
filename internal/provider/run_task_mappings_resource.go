package provider

import (
	"context"
	"fmt"

	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	sdkRunTask "github.com/control-monkey/controlmonkey-sdk-go/services/run_task"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	runTaskMapping "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/run_task_mappings"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
var _ resource.Resource = &RunTaskMappingResource{}

func NewRunTaskMappingResource() resource.Resource {
	return &RunTaskMappingResource{}
}

type RunTaskMappingResource struct {
	client *ControlMonkeyAPIClient
}

func (r *RunTaskMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_run_task_mappings"
}

func (r *RunTaskMappingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys run task mappings. For more information: [ControlMonkey Documentation](https://docs.controlmonkey.io/main-concepts/run-tasks)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"run_task_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the run task.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"targets": commons.WithNullSetDefault(schema.SetNestedAttribute{
				MarkdownDescription: "List of targets",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"target_id": schema.StringAttribute{
							MarkdownDescription: fmt.Sprintf("The unique ID corresponds to the `target_type` in the mapping."+
								" Use `%s` with `target_type` `%s` to map the run task to all namespaces in the organization.",
								cmTypes.AllTargetsIdentifier, cmTypes.NamespaceTargetType),
							Required: true,
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
							MarkdownDescription: fmt.Sprintf("Specifies the level of enforcement for the run task on the target. Allowed values: %s."+
								" When set to `softMandatory`, a run task failure triggers an approval requirement before applying changes."+
								" When set to `hardMandatory`, changes cannot be applied until the run task is successful.", helpers.EnumForDocs(cmTypes.EnforcementLevelTypes)),
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(cmTypes.EnforcementLevelTypes...),
							},
						},
						"stage": schema.StringAttribute{
							MarkdownDescription: "The stage in which the run task will execute. Find supported types [here](https://docs.controlmonkey.io/controlmonkey-api/api-enumerations#run-task-stage)",
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
func (r *RunTaskMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RunTaskMappingResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data runTaskMapping.ResourceModel

	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		return
	}

	if len(data.Targets) > 0 {
		identifiers := interfaces.GetIdentifiers(data.Targets)

		if helpers.IsUnique(identifiers) == false {
			duplicates := helpers.FindDuplicates(identifiers, false)
			for _, d := range duplicates {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("Target '%s' appears more than once", runTaskMapping.CleanIdentifier(d)))
			}
		}

		for _, t := range data.Targets {
			if targetType := t.TargetType; helpers.IsKnown(targetType) {
				if t.TargetId.ValueString() == cmTypes.AllTargetsIdentifier && targetType.ValueString() != cmTypes.NamespaceTargetType {
					resp.Diagnostics.AddError(validationError, fmt.Sprintf("Only target_type '%s' can have target_id value '%s'",
						cmTypes.NamespaceTargetType, cmTypes.AllTargetsIdentifier))
				}
			}
		}
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *RunTaskMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state runTaskMapping.ResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.runTask.ListRunTaskMappings(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(runTaskNotFoundError, fmt.Sprintf("Run task '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read mappings for run task '%s'", id), err.Error())
		return
	}

	runTaskMapping.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *RunTaskMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan runTaskMapping.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := runTaskMapping.Merge(&plan, nil, commons.CreateConverter)
	runTaskId := plan.RunTaskId

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, runTaskId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = runTaskId

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RunTaskMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan runTaskMapping.ResourceModel
	var state runTaskMapping.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := runTaskMapping.Merge(&plan, &state, commons.UpdateMerger)
	runTaskId := plan.RunTaskId

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, runTaskId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.updateEntities(ctx, mergeResult.EntitiesToUpdate, runTaskId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, runTaskId.ValueString())
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

func (r *RunTaskMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state runTaskMapping.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := runTaskMapping.Merge(nil, &state, commons.DeleteMerger)
	runTaskId := state.RunTaskId

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, runTaskId.ValueString())
	resp.Diagnostics.Append(diags...)
}

func (r *RunTaskMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *RunTaskMappingResource) createEntities(ctx context.Context, entitiesToCreate []*sdkRunTask.RunTaskMapping, runTaskId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Mapping %d targets to run task '%s'.", len(entitiesToCreate), runTaskId))

	for _, e := range entitiesToCreate {
		_, err := r.client.Client.runTask.CreateRunTaskMapping(ctx, e)

		if err != nil {
			targetId := *e.TargetId
			targetType := *e.TargetType
			if commons.IsAlreadyExistResponseError(err) {
				tflog.Info(ctx, fmt.Sprintf("Target '%s' of type '%s' is already mapped to run task '%s'. No operation was made.", targetId, targetType, runTaskId))
			} else if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to create map between target '%s' of type '%s' and run task '%s'. Error: %s", targetId, targetType, runTaskId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to create map between target '%s' of type '%s' and run task '%s'", targetId, targetType, runTaskId),
						err.Error()),
				}
			}
		}
	}

	return retVal
}

func (r *RunTaskMappingResource) updateEntities(ctx context.Context, entitiesToUpdate []*sdkRunTask.RunTaskMapping, runTaskId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Updating %d target mappings to run task '%s'.", len(entitiesToUpdate), runTaskId))

	for _, e := range entitiesToUpdate {
		_, err := r.client.Client.runTask.UpdateRunTaskMapping(ctx, e)

		if err != nil {
			targetId := *e.TargetId
			targetType := *e.TargetType
			if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to update map between target '%s' of type '%s' and run task '%s'. Error: %s", targetId, targetType, runTaskId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to update map between target '%s' of type '%s' and run task '%s'", targetId, targetType, runTaskId),
						err.Error()),
				}
			}
		}
	}

	return retVal
}

func (r *RunTaskMappingResource) deleteEntities(ctx context.Context, entitiesToDelete []*sdkRunTask.RunTaskMapping, runTaskId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Removing %d target mappings from run task '%s'.", len(entitiesToDelete), runTaskId))

	for _, e := range entitiesToDelete {
		_, err := r.client.Client.runTask.DeleteRunTaskMapping(ctx, e)

		if err != nil {
			targetId := *e.TargetId
			targetType := *e.TargetType
			if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to delete map between target '%s' of type '%s' and run task '%s'. Error: %s", targetId, targetType, runTaskId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to delete map between target '%s' of type '%s' and run task '%s'", targetId, targetType, runTaskId),
						err.Error()),
				}
			}
		}
	}

	return retVal
}
