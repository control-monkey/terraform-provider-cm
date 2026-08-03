package provider

import (
	"context"
	"fmt"

	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	sdkCustomFlow "github.com/control-monkey/controlmonkey-sdk-go/services/custom_flow"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	customFlowMapping "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/custom_flow_mappings"
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
var _ resource.Resource = &CustomFlowMappingResource{}

func NewCustomFlowMappingResource() resource.Resource {
	return &CustomFlowMappingResource{}
}

type CustomFlowMappingResource struct {
	client *ControlMonkeyAPIClient
}

func (r *CustomFlowMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_flow_mappings"
}

func (r *CustomFlowMappingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates and destroys custom flow mappings. For more information: [ControlMonkey Documentation](https://docs.controlmonkey.io/main-concepts/runs/custom-flows)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_flow_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the custom flow.",
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
								" Use `%s` with `target_type` `%s` to map the custom flow to all namespaces in the organization.",
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
					},
				},
			}),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *CustomFlowMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomFlowMappingResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data customFlowMapping.ResourceModel

	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		return
	}

	if len(data.Targets) > 0 {
		identifiers := interfaces.GetIdentifiers(data.Targets)

		if helpers.IsUnique(identifiers) == false {
			duplicates := helpers.FindDuplicates(identifiers, false)
			for _, d := range duplicates {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("Target '%s' appears more than once", customFlowMapping.CleanIdentifier(d)))
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
func (r *CustomFlowMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state customFlowMapping.ResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.customFlow.ListCustomFlowMappings(ctx)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(customFlowNotFoundError, fmt.Sprintf("Custom flow '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read mappings for custom flow '%s'", id), err.Error())
		return
	}

	customFlowMapping.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *CustomFlowMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan customFlowMapping.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := customFlowMapping.Merge(&plan, nil, commons.CreateConverter)
	customFlowId := plan.CustomFlowId

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, customFlowId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = customFlowId

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update handles target changes. The API has no update endpoint, so a changed target is a
// delete followed by a create.
func (r *CustomFlowMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan customFlowMapping.ResourceModel
	var state customFlowMapping.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := customFlowMapping.Merge(&plan, &state, commons.UpdateMerger)
	customFlowId := plan.CustomFlowId

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, customFlowId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, customFlowId.ValueString())
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

func (r *CustomFlowMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state customFlowMapping.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := customFlowMapping.Merge(nil, &state, commons.DeleteMerger)
	customFlowId := state.CustomFlowId

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, customFlowId.ValueString())
	resp.Diagnostics.Append(diags...)
}

func (r *CustomFlowMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CustomFlowMappingResource) createEntities(ctx context.Context, entitiesToCreate []*sdkCustomFlow.CustomFlowMapping, customFlowId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Mapping %d targets to custom flow '%s'.", len(entitiesToCreate), customFlowId))

	for _, e := range entitiesToCreate {
		_, err := r.client.Client.customFlow.CreateCustomFlowMapping(ctx, e)

		if err != nil {
			targetId := *e.TargetId
			targetType := *e.TargetType
			if commons.IsAlreadyExistResponseError(err) {
				tflog.Info(ctx, fmt.Sprintf("Target '%s' of type '%s' is already mapped to custom flow '%s'. No operation was made.", targetId, targetType, customFlowId))
			} else if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to create map between target '%s' of type '%s' and custom flow '%s'. Error: %s", targetId, targetType, customFlowId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to create map between target '%s' of type '%s' and custom flow '%s'", targetId, targetType, customFlowId),
						err.Error()),
				}
			}
		}
	}

	return retVal
}

func (r *CustomFlowMappingResource) deleteEntities(ctx context.Context, entitiesToDelete []*sdkCustomFlow.CustomFlowMapping, customFlowId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Removing %d target mappings from custom flow '%s'.", len(entitiesToDelete), customFlowId))

	for _, e := range entitiesToDelete {
		_, err := r.client.Client.customFlow.DeleteCustomFlowMapping(ctx, e)

		if err != nil {
			targetId := *e.TargetId
			targetType := *e.TargetType
			if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to delete map between target '%s' of type '%s' and custom flow '%s'. Error: %s", targetId, targetType, customFlowId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to delete map between target '%s' of type '%s' and custom flow '%s'", targetId, targetType, customFlowId),
						err.Error()),
				}
			}
		}
	}

	return retVal
}
