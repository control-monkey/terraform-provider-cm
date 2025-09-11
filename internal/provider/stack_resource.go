package provider

import (
	"context"
	"fmt"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/cross_schema"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/stack"
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
var _ resource.Resource = &StackResource{}

func NewStackResource() resource.Resource {
	return &StackResource{}
}

type StackResource struct {
	client *ControlMonkeyAPIClient
}

func (r *StackResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stack"
}

func (r *StackResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys stacks. For more information: [ControlMonkey Documentation](https://docs.controlmonkey.io/main-concepts/stack)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the stack.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"iac_type": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("IaC type of the stack. Allowed values: %s.", helpers.EnumForDocs(cmTypes.IacTypes)),
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(cmTypes.IacTypes...),
				},
			},
			"namespace_id": schema.StringAttribute{
				MarkdownDescription: "The namespace ID where the stack is located.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the stack.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the stack.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.NoneOf(""),
				},
			},
			"deployment_behavior":        cross_schema.StackDeploymentBehaviorSchema,
			"deployment_approval_policy": cross_schema.StackDeploymentApprovalPolicySchema,
			"vcs_info": schema.SingleNestedAttribute{
				MarkdownDescription: "The configuration of the version control to which the stack is attached.",
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
						MarkdownDescription: "The path to a chosen directory from the root. Default path is root directory",
						Optional:            true,
					},
					"branch": schema.StringAttribute{
						MarkdownDescription: "The branch that should trigger plan/deployment for the stack. When no branch is given, the default branch of the repository is chosen.",
						Optional:            true,
					},
				},
			},
			"run_trigger": cross_schema.RunTriggerSchema,
			"iac_config":  cross_schema.IacConfigSchema,
			"policy": schema.SingleNestedAttribute{
				MarkdownDescription: "The policy of the stack.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"ttl_config": schema.SingleNestedAttribute{
						MarkdownDescription: "The time to live config of the stack policy.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"ttl": schema.SingleNestedAttribute{
								Required: true,
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										MarkdownDescription: fmt.Sprintf("The type of the ttl. Allowed values: %s.", helpers.EnumForDocs(cmTypes.TtlTypes)),
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf(cmTypes.TtlTypes...),
										},
									},
									"value": schema.Int64Attribute{
										MarkdownDescription: "The value that corresponds the type",
										Required:            true,
									},
								},
							},
						},
					},
				},
			},
			"runner_config": cross_schema.StackRunnerConfigSchema,
			"auto_sync":     cross_schema.AutoSyncSchema,
			"capabilities": schema.SingleNestedAttribute{
				MarkdownDescription: "List of capabilities enabled for the stack.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"deploy_on_push": schema.SingleNestedAttribute{
						MarkdownDescription: "When enabled, a deployment will be automatically triggered when changes are pushed to the repository that are relevant to the stack.",
						Optional:            true,
						Attributes:          stackCapabilityConfigSchema(),
					},
					"plan_on_pr": schema.SingleNestedAttribute{
						MarkdownDescription: "When enabled, a plan will be automatically triggered when a Pull Request is created or updated with changes relevant to the stack.",
						Optional:            true,
						Attributes:          stackCapabilityConfigSchema(),
					},
					"drift_detection": schema.SingleNestedAttribute{
						MarkdownDescription: "When enabled, ControlMonkey will frequently check for drifts in your stack configuration.",
						Optional:            true,
						Attributes:          stackCapabilityConfigSchema(),
					},
				},
			},
		},
	}
}

func stackCapabilityConfigSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"status": schema.StringAttribute{
			MarkdownDescription: fmt.Sprint("Whether the capability is enabled or disabled. Allowed values: [enabled, disabled]."),
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("enabled", "disabled"),
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *StackResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StackResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data stack.ResourceModel

	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		return
	}

	runnerConfig := data.RunnerConfig

	if runnerConfig != nil {
		mode := runnerConfig.Mode

		if helpers.IsKnown(mode) {
			modeValue := mode.ValueString()

			if modeValue == cmTypes.Managed && runnerConfig.Groups.IsNull() == false {
				resp.Diagnostics.AddError(
					validationError, fmt.Sprintf("runner_config.mode with type '%s' cannot have runner_config.groups", cmTypes.Managed),
				)
			} else if modeValue == cmTypes.SelfHosted && helpers.IsKnown(runnerConfig.Groups) {
				if len(runnerConfig.Groups.Elements()) == 0 {
					resp.Diagnostics.AddError(
						validationError, fmt.Sprintf("runner_config.mode with type '%s' requires runner_config.groups to be not empty", cmTypes.SelfHosted),
					)
				} else if helpers.DoesTfListContainsEmptyValue(runnerConfig.Groups) {
					resp.Diagnostics.AddError(
						validationError, "Found empty string in runner_config.groups",
					)
				} else if !helpers.IsTfStringSliceUnique(runnerConfig.Groups) {
					resp.Diagnostics.AddError(
						validationError, "Found duplicate in runner_config.groups",
					)
				}
			}
		}
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *StackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get current state
	var state stack.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.stack.ReadStack(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read stack %s", id), fmt.Sprintf("%s", err))
		return
	}

	stack.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *StackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan stack.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := stack.Converter(&plan, nil, commons.CreateConverter)

	res, err := r.client.Client.stack.CreateStack(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Stack creation failed",
			fmt.Sprintf("failed to create stack, error: %s", err.Error()),
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

func (r *StackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan stack.ResourceModel
	var state stack.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body, _ := stack.Converter(&plan, &state, commons.UpdateConverter)

	_, err := r.client.Client.stack.UpdateStack(ctx, id, body)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.Diagnostics.AddError(resourceNotFoundError, fmt.Sprintf("Stack '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(
			resourceUpdateFailedError,
			fmt.Sprintf("failed to update stack %s, error: %s", id, err),
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

func (r *StackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state stack.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	_, err := r.client.Client.stack.DeleteStack(ctx, id)

	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Stack deletion failed",
			fmt.Sprintf("Failed to delete stack %s, error: %s", id, err),
		)
		return
	}
}

func (r *StackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
