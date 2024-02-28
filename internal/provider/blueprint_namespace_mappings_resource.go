package provider

import (
	"context"
	"fmt"
	"github.com/control-monkey/controlmonkey-sdk-go/services/blueprint"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	blueprintNamespaces "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/blueprint_namespace_mappings"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
var _ resource.Resource = &BlueprintNamespaceMappingsResource{}

func NewBlueprintNamespaceMappingsResource() resource.Resource {
	return &BlueprintNamespaceMappingsResource{}
}

type BlueprintNamespaceMappingsResource struct {
	client *ControlMonkeyAPIClient
}

func (r *BlueprintNamespaceMappingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blueprint_namespace_mappings"
}

func (r *BlueprintNamespaceMappingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys mappings from blueprint to namespaces.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of this resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the blueprint.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"namespaces": schema.SetNestedAttribute{
				MarkdownDescription: "A list of namespaces to which the blueprint is mapped.",
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
						"namespace_id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the namespace.",
							Required:            true,
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
func (r *BlueprintNamespaceMappingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BlueprintNamespaceMappingsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data blueprintNamespaces.ResourceModel

	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		return
	}

	if len(data.Namespaces) > 0 {
		identifiers := interfaces.GetIdentifiers(data.Namespaces)

		if helpers.IsUnique(identifiers) == false {
			duplicates := helpers.FindDuplicates(identifiers, false)
			for _, d := range duplicates {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("Namespace with id '%s' appears more than once", d))
			}
		}
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *BlueprintNamespaceMappingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get current state
	var state blueprintNamespaces.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.blueprint.ListBlueprintNamespaceMappings(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.Diagnostics.AddWarning(blueprintNotFoundError, fmt.Sprintf("Blueprint '%s' not found", id))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read blueprint-namespace mappings for blueprint '%s'", id), err.Error())
		return
	}

	blueprintNamespaces.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *BlueprintNamespaceMappingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan blueprintNamespaces.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := blueprintNamespaces.Merge(&plan, nil, commons.CreateMerger)
	blueprintId := plan.BlueprintId

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, blueprintId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = blueprintId

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *BlueprintNamespaceMappingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan blueprintNamespaces.ResourceModel
	var state blueprintNamespaces.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := blueprintNamespaces.Merge(&plan, &state, commons.UpdateMerger)

	diags = r.createEntities(ctx, mergeResult.EntitiesToCreate, plan.BlueprintId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, plan.BlueprintId.ValueString())
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

func (r *BlueprintNamespaceMappingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state blueprintNamespaces.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := blueprintNamespaces.Merge(nil, &state, commons.DeleteMerger)

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, state.BlueprintId.ValueString())
	resp.Diagnostics.Append(diags...)
}

func (r *BlueprintNamespaceMappingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

//region Private Methods

func (r *BlueprintNamespaceMappingsResource) createEntities(ctx context.Context, entitiesToCreate []*blueprint.BlueprintNamespaceMapping, blueprintId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Mapping %d namespaces to blueprint '%s'.", len(entitiesToCreate), blueprintId))

	for _, e := range entitiesToCreate {
		_, err := r.client.Client.blueprint.CreateBlueprintNamespaceMapping(ctx, e)

		if err != nil {
			namespaceId := *e.NamespaceId
			if commons.IsAlreadyExistResponseError(err) {
				tflog.Info(ctx, fmt.Sprintf("Namespace '%s' is already mapped to blueprint '%s'. No operation was made.", namespaceId, blueprintId))
			} else if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to map namespace '%s' to blueprint '%s'. Error: %s",
						namespaceId, blueprintId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to map namespace '%s' to blueprint '%s'", namespaceId, blueprintId), err.Error()),
				}
			}
		}
	}

	return retVal
}

func (r *BlueprintNamespaceMappingsResource) deleteEntities(ctx context.Context, entitiesToDelete []*blueprint.BlueprintNamespaceMapping, blueprintId string) diag.Diagnostics {
	var retVal diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Removing %d namespace mappings from blueprint '%s'.", len(entitiesToDelete), blueprintId))

	for _, e := range entitiesToDelete {
		_, err := r.client.Client.blueprint.DeleteBlueprintNamespaceMapping(ctx, e)

		if err != nil {
			namespaceId := *e.NamespaceId
			if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to delete mapping between namespace '%s' and blueprint '%s'. Error: %s", namespaceId, blueprintId, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to delete mapping between namespace '%s' and blueprint '%s'", namespaceId, blueprintId),
						err.Error()),
				}
			}
		}
	}

	return retVal
}

//endregion
