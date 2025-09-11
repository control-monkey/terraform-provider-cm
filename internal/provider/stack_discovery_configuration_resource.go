package provider

import (
	"context"
	"fmt"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/cross_schema"
	tfStackDiscoveryConfiguration "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/stack_discovery_configuration"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &StackDiscoveryConfigurationResource{}

func NewStackDiscoveryConfigurationResource() resource.Resource {
	return &StackDiscoveryConfigurationResource{}
}

type StackDiscoveryConfigurationResource struct {
	client *ControlMonkeyAPIClient
}

func (r *StackDiscoveryConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stack_discovery_configuration"
}

func (r *StackDiscoveryConfigurationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys stack discovery configurations. For more information: [ControlMonkey Documentation](https://docs.controlmonkey.io/main-concepts/stack/stack-auto-discovery)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the stack discovery configuration.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the stack discovery configuration.",
				Required:            true,
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"namespace_id": schema.StringAttribute{
				MarkdownDescription: "The namespace ID where the stack discovery configuration is located.",
				Required:            true,
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the stack discovery configuration.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.NoneOf(""),
				},
			},
			"vcs_patterns": schema.ListNestedAttribute{
				MarkdownDescription: "The VCS patterns configuration for stack discovery.",
				Required:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"provider_id": schema.StringAttribute{
							MarkdownDescription: "The ControlMonkey unique ID of the connected version control system.",
							Required:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
						"repo_name": schema.StringAttribute{
							MarkdownDescription: "The name of the version control repository.",
							Required:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
						"path_patterns": schema.ListAttribute{
							MarkdownDescription: "List of path patterns to include for stack discovery.",
							ElementType:         types.StringType,
							Required:            true,
							Validators:          commons.ValidateUniqueNotEmptyListWithNoBlankValues(),
						},
						"exclude_path_patterns": schema.ListAttribute{
							MarkdownDescription: "List of path patterns to exclude from stack discovery.",
							ElementType:         types.StringType,
							Optional:            true,
							Validators:          commons.ValidateUniqueNotEmptyListWithNoBlankValues(),
						},
						"branch": schema.StringAttribute{
							MarkdownDescription: "The branch to monitor for stack discovery.",
							Required:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
					},
				},
			},
			"stack_config": schema.SingleNestedAttribute{
				MarkdownDescription: "The stack configuration template for discovered stacks.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"iac_type": schema.StringAttribute{
						MarkdownDescription: fmt.Sprintf("IaC type of the stack. Allowed values: %s.", helpers.EnumForDocs(cmTypes.IacTypes)),
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(cmTypes.IacTypes...),
						},
					},
					"deployment_behavior":        cross_schema.StackDeploymentBehaviorSchema,
					"deployment_approval_policy": cross_schema.StackDeploymentApprovalPolicySchema,
					"run_trigger":                cross_schema.RunTriggerSchema,
					"iac_config":                 cross_schema.IacConfigSchema,
					"runner_config":              cross_schema.StackRunnerConfigSchema,
					"auto_sync":                  cross_schema.AutoSyncSchema,
				},
			},
		},
	}
}

func (r *StackDiscoveryConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StackDiscoveryConfigurationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data tfStackDiscoveryConfiguration.ResourceModel

	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		return
	}

	if data.StackConfig != nil && data.StackConfig.RunnerConfig != nil {
		runnerConfig := data.StackConfig.RunnerConfig
		mode := runnerConfig.Mode

		if helpers.IsKnown(mode) {
			modeValue := mode.ValueString()

			if modeValue == cmTypes.Managed && runnerConfig.Groups.IsNull() == false {
				resp.Diagnostics.AddError(
					validationError, fmt.Sprintf("stack_config.runner_config.mode with type '%s' cannot have stack_config.runner_config.groups", cmTypes.Managed),
				)
			} else if modeValue == cmTypes.SelfHosted && helpers.IsKnown(runnerConfig.Groups) {
				if len(runnerConfig.Groups.Elements()) == 0 {
					resp.Diagnostics.AddError(
						validationError, fmt.Sprintf("stack_config.runner_config.mode with type '%s' requires stack_config.runner_config.groups to be not empty", cmTypes.SelfHosted),
					)
				} else if helpers.DoesTfListContainsEmptyValue(runnerConfig.Groups) {
					resp.Diagnostics.AddError(
						validationError, "Found empty string in stack_config.runner_config.groups",
					)
				} else if !helpers.IsTfStringSliceUnique(runnerConfig.Groups) {
					resp.Diagnostics.AddError(
						validationError, "Found duplicate in stack_config.runner_config.groups",
					)
				}
			}
		}
	}

	// Validate path_patterns elements are not blank
	if data.VcsPatterns != nil {
		for i, vcsPattern := range data.VcsPatterns {
			if helpers.IsKnown(vcsPattern.PathPatterns) {
				if helpers.DoesTfListContainsEmptyValue(vcsPattern.PathPatterns) {
					resp.Diagnostics.AddError(
						validationError, fmt.Sprintf("Found empty string in vcs_patterns[%d].path_patterns", i),
					)
				}
			}
			if helpers.IsKnown(vcsPattern.ExcludePathPatterns) {
				if helpers.DoesTfListContainsEmptyValue(vcsPattern.ExcludePathPatterns) {
					resp.Diagnostics.AddError(
						validationError, fmt.Sprintf("Found empty string in vcs_patterns[%d].exclude_path_patterns", i),
					)
				}
			}
		}
	}
}

func (r *StackDiscoveryConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tfStackDiscoveryConfiguration.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.stackDiscoveryConfiguration.ReadStackDiscoveryConfiguration(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read stack discovery configuration %s", id), fmt.Sprintf("%s", err))
		return
	}

	tfStackDiscoveryConfiguration.UpdateStateAfterRead(res, &state)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *StackDiscoveryConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tfStackDiscoveryConfiguration.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := tfStackDiscoveryConfiguration.Converter(&plan, nil, commons.CreateConverter)

	res, err := r.client.Client.stackDiscoveryConfiguration.CreateStackDiscoveryConfiguration(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Stack discovery configuration creation failed",
			fmt.Sprintf("failed to create stack discovery configuration, error: %s", err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(controlmonkey.StringValue(res.ID))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *StackDiscoveryConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tfStackDiscoveryConfiguration.ResourceModel
	var state tfStackDiscoveryConfiguration.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body, _ := tfStackDiscoveryConfiguration.Converter(&plan, &state, commons.UpdateConverter)

	_, err := r.client.Client.stackDiscoveryConfiguration.UpdateStackDiscoveryConfiguration(ctx, id, body)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.Diagnostics.AddError(resourceNotFoundError, fmt.Sprintf("Stack discovery configuration '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(
			resourceUpdateFailedError,
			fmt.Sprintf("failed to update stack discovery configuration %s, error: %s", id, err),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *StackDiscoveryConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tfStackDiscoveryConfiguration.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	_, err := r.client.Client.stackDiscoveryConfiguration.DeleteStackDiscoveryConfiguration(ctx, id)

	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Stack discovery configuration deletion failed",
			fmt.Sprintf("Failed to delete stack discovery configuration %s, error: %s", id, err),
		)
		return
	}
}

func (r *StackDiscoveryConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
