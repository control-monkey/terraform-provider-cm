package provider

import (
	"context"
	"fmt"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/cross_schema"
	tfBlueprint "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/blueprint"
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
var _ resource.Resource = &BlueprintResource{}

func NewBlueprintResource() resource.Resource {
	return &BlueprintResource{}
}

type BlueprintResource struct {
	client *ControlMonkeyAPIClient
}

func (r *BlueprintResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blueprint"
}

func (r *BlueprintResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys blueprints. For more information: [ControlMonkey Documentation](https://docs.controlmonkey.io/main-concepts/self-service-templates/persistent-template)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the blueprint.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the blueprint.",
				Required:            true,
				Validators:          []validator.String{cm_stringvalidators.NotBlank()},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the blueprint.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.NoneOf(""),
				},
			},
			"blueprint_vcs_info": schema.SingleNestedAttribute{
				MarkdownDescription: "Configuration details for the version control system storing the blueprint.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"provider_id": schema.StringAttribute{
						MarkdownDescription: "The ControlMonkey unique ID of the connected version control system.",
						Required:            true,
					},
					"repo_name": schema.StringAttribute{
						MarkdownDescription: "The name of the version control repository.",
						Required:            true,
						Validators:          []validator.String{cm_stringvalidators.NotBlank()},
					},
					"path": schema.StringAttribute{
						MarkdownDescription: "The relative path to the directory containing the blueprint files, starting from the root of the repository.",
						Required:            true,
						Validators:          []validator.String{cm_stringvalidators.NotBlank()},
					},
					"branch": schema.StringAttribute{
						MarkdownDescription: "The branch in which the blueprint is located. When no branch is given, the default branch of the repository is chosen.",
						Optional:            true,
						Validators:          []validator.String{cm_stringvalidators.NotBlank()},
					},
				},
			},
			"stack_configuration": schema.SingleNestedAttribute{
				MarkdownDescription: "The configuration for creating new persistent stacks from the blueprint.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"name_pattern": schema.StringAttribute{
						MarkdownDescription: "A pattern used to name persistent stacks created from the blueprint. The pattern must include at least one dynamic substitute parameter (e.g., `{region}-{service}`).",
						Required:            true,
					},
					"iac_type": schema.StringAttribute{
						MarkdownDescription: fmt.Sprintf("IaC type of the template. Allowed values: %s.", helpers.EnumForDocs(cmTypes.IacTypes)),
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(cmTypes.IacTypes...),
						},
					},
					"vcs_info_with_patterns": schema.SingleNestedAttribute{
						MarkdownDescription: "Configuration details for the version control system where the stack files generated from the blueprint will be stored.",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"provider_id": schema.StringAttribute{
								MarkdownDescription: "The ControlMonkey unique ID of the connected version control system.",
								Required:            true,
							},
							"repo_name": schema.StringAttribute{
								MarkdownDescription: "The name of the version control repository.",
								Required:            true,
								Validators:          []validator.String{cm_stringvalidators.NotBlank()},
							},
							"path_pattern": schema.StringAttribute{
								MarkdownDescription: "A pattern to a new path in the repository to which new persistent stack files created from the blueprint will be pushed. This field requires at least one substitute parameter.",
								Required:            true,
								Validators:          []validator.String{cm_stringvalidators.NotBlank()},
							},
							"branch_pattern": schema.StringAttribute{
								MarkdownDescription: "The target branch for new pull requests containing the new stack files. Substitute parameters (e.g., `{branch}-{env}`) are supported.",
								Optional:            true,
							},
						},
					},
					"deployment_approval_policy": cross_schema.StackDeploymentApprovalPolicySchema,
					"run_trigger":                cross_schema.RunTriggerSchema,
					"iac_config":                 cross_schema.IacConfigSchema,
					"auto_sync":                  cross_schema.AutoSyncSchema,
				},
			},
			"substitute_parameters": schema.ListNestedAttribute{
				MarkdownDescription: "Define dynamic placeholders (`{parameter_name}`) used in patterns (e.g., `name_pattern`, `path_pattern`) or Terraform files. Users will supply values for these parameters when launching stacks.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							MarkdownDescription: "The key for the substitute parameter excluding the curly braces. For example, if the Terraform file contains `{replace-me}`, the key should be `replace-me`.",
							Required:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "A description of the parameter. Users launching stacks from this blueprint will reference this description to assign values. Providing a clear, meaningful description is highly recommended.",
							Required:            true,
						},
						"value_conditions": cross_schema.ValueConditionsSchema,
					},
				},
			},
			"skip_plan_on_stack_initialization": schema.BoolAttribute{
				MarkdownDescription: "If enabled (`true`), an automatic plan will not be triggered on the initial pull request.",
				Optional:            true,
			},
			"auto_approve_apply_on_initialization": schema.BoolAttribute{
				MarkdownDescription: "If enabled (`true`), the stack’s initial deployment will automatically apply changes after the pull request is merged, bypassing manual approval.",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *BlueprintResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *BlueprintResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state tfBlueprint.ResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.blueprint.ReadBlueprint(ctx, id)

	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(resourceNotFoundError, fmt.Sprintf("Blueprint '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read blueprint '%s'", id), err.Error())
		return
	}

	tfBlueprint.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *BlueprintResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan tfBlueprint.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := tfBlueprint.Converter(&plan, nil, commons.CreateConverter)

	res, err := r.client.Client.blueprint.CreateBlueprint(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceCreationFailedError,
			fmt.Sprintf("failed to create blueprint, error: %s", err.Error()),
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

func (r *BlueprintResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan tfBlueprint.ResourceModel
	var state tfBlueprint.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body, _ := tfBlueprint.Converter(&plan, &state, commons.UpdateConverter)

	_, err := r.client.Client.blueprint.UpdateBlueprint(ctx, id, body)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.Diagnostics.AddError(resourceNotFoundError, fmt.Sprintf("Blueprint '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(
			resourceUpdateFailedError,
			fmt.Sprintf("failed to update blueprint %s, error: %s", id, err),
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

func (r *BlueprintResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state tfBlueprint.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	_, err := r.client.Client.blueprint.DeleteBlueprint(ctx, id)

	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			resourceDeletionFailedError,
			fmt.Sprintf("Failed to delete blueprint %s, error: %s", id, err),
		)
		return
	}
}

func (r *BlueprintResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
