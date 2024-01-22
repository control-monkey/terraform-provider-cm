package provider

import (
	"context"
	"fmt"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/namespace"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
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
var _ resource.Resource = &NamespaceResource{}

func NewNamespaceResource() resource.Resource {
	return &NamespaceResource{}
}

type NamespaceResource struct {
	client *ControlMonkeyAPIClient
}

func (r *NamespaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}

func (r *NamespaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys namespaces.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the namespace.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the namespace.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the namespace.",
				Optional:            true,
			},
			"external_credentials": schema.ListNestedAttribute{
				MarkdownDescription: "List of cloud credentials attached to the namespace.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: fmt.Sprintf("The type of the credentials. Allowed values: %s.", helpers.EnumForDocs(cmTypes.ExternalCredentialTypes)),
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(cmTypes.ExternalCredentialTypes...),
							},
						},
						"external_credentials_id": schema.StringAttribute{
							MarkdownDescription: "The ControlMonkey unique ID of the credentials.",
							Required:            true,
						},
						"aws_profile_name": schema.StringAttribute{
							MarkdownDescription: "Profile name for AWS credentials.",
							Optional:            true,
						},
					},
				},
				Validators: []validator.List{listvalidator.SizeAtLeast(1)},
			},
			"policy": schema.SingleNestedAttribute{
				MarkdownDescription: "The policy of the namespace.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"ttl_config": schema.SingleNestedAttribute{
						MarkdownDescription: "The time to live config of the namespace policy regarding to its stacks.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"max_ttl": schema.SingleNestedAttribute{
								MarkdownDescription: "The max time to live for new stacks in the namespace.",
								Required:            true,
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
							"default_ttl": schema.SingleNestedAttribute{
								MarkdownDescription: "The default time to live for new stacks in the namespace.",
								Required:            true,
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
			"iac_config": schema.SingleNestedAttribute{
				MarkdownDescription: "IaC configuration of the namespace. If not overridden, this becomes the default for its stacks.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"terraform_version": schema.StringAttribute{
						MarkdownDescription: "the Terraform version that will be used for terraform operations.",
						Optional:            true,
					},
					"terragrunt_version": schema.StringAttribute{
						MarkdownDescription: "the Terragrunt version that will be used for terragrunt operations.",
						Optional:            true,
					},
					"opentofu_version": schema.StringAttribute{
						MarkdownDescription: "the OpenTofu version that will be used for tofu operations.",
						Optional:            true,
					},
				},
			},
			"runner_config": schema.SingleNestedAttribute{
				MarkdownDescription: "Configure the runner settings to specify whether ControlMonkey manages the runner or it is self-hosted.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"mode": schema.StringAttribute{
						MarkdownDescription: fmt.Sprintf("The runner mode. Allowed values: %s.", helpers.EnumForDocs(cmTypes.RunnerConfigModeTypes)),
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(cmTypes.RunnerConfigModeTypes...),
						},
					},
					"groups": schema.ListAttribute{
						MarkdownDescription: fmt.Sprintf("In case that `mode` is `%s`, groups must contain at least one runners group. If `mode` is `%s`, this field must not be configures.", cmTypes.SelfHosted, cmTypes.Managed),
						ElementType:         types.StringType,
						Optional:            true,
						Validators:          commons.ValidateUniqueNotEmptyListWithNoBlankValues(),
					},
					"is_overridable": schema.BoolAttribute{
						MarkdownDescription: "Determine if stacks within the namespace can override the runner_config.",
						Required:            true,
					},
				},
			},
			"deployment_approval_policy": schema.SingleNestedAttribute{
				MarkdownDescription: "Set up requirements to approve a deployment",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"rules": schema.ListNestedAttribute{
						MarkdownDescription: "Set up rules for approving deployment processes. At least one rule should be configured",
						Required:            true,
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									MarkdownDescription: fmt.Sprintf("The type of the rule. Allowed values: %s.", helpers.EnumForDocs(cmTypes.DeploymentApprovalPolicyRuleTypes)),
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf(cmTypes.DeploymentApprovalPolicyRuleTypes...),
									},
								},
							},
						},
					},
					"override_behavior": schema.StringAttribute{
						MarkdownDescription: fmt.Sprintf("Decide whether stacks can override this configuration. Allowed values: %s.", helpers.EnumForDocs(cmTypes.OverrideBehaviorTypes)),
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(cmTypes.OverrideBehaviorTypes...),
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *NamespaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NamespaceResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data namespace.ResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	externalCredentials := data.ExternalCredentials

	if externalCredentials != nil {
		for _, credentials := range externalCredentials {
			credentialsType := credentials.Type
			profileName := credentials.AwsProfileName

			if credentialsType.IsNull() == false && credentialsType.ValueString() != cmTypes.AwsAssumeRole && profileName.IsNull() == false {
				resp.Diagnostics.AddError(
					"Validation Error",
					fmt.Sprintf("external_credentials cannot have aws_profile_name configured for non AWS provider."),
				)
			}
		}
	}

	runnerConfig := data.RunnerConfig

	if runnerConfig != nil {
		mode := runnerConfig.Mode

		if !mode.IsNull() {
			modeValue := mode.ValueString()

			if modeValue == cmTypes.Managed && runnerConfig.Groups != nil {
				resp.Diagnostics.AddError(
					"Validation Error",
					fmt.Sprintf("runner_config.mode with type '%s' cannot have runner_config.groups", cmTypes.Managed),
				)
			} else if modeValue == cmTypes.SelfHosted {
				if len(runnerConfig.Groups) == 0 {
					resp.Diagnostics.AddError(
						"Validation Error",
						fmt.Sprintf("runner_config.mode with type '%s' requires runner_config.groups to be not empty", cmTypes.SelfHosted),
					)
				} else if helpers.DoesTfStringSliceContainEmptyValue(runnerConfig.Groups) {
					resp.Diagnostics.AddError(
						"Validation Error",
						"Found empty string in runner_config.groups",
					)
				} else if !helpers.IsTfStringSliceUnique(runnerConfig.Groups) {
					resp.Diagnostics.AddError(
						"Validation Error",
						"Found duplicate in runner_config.groups",
					)
				}
			}
		}
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *NamespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get current state
	var state namespace.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.namespace.ReadNamespace(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read namespace %s", id), fmt.Sprintf("%s", err))
		return
	}

	namespace.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *NamespaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan namespace.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := namespace.Converter(&plan, nil, commons.CreateConverter)

	res, err := r.client.Client.namespace.CreateNamespace(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Namespace creation failed",
			fmt.Sprintf("failed to create namespace, error: %s", err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(controlmonkey.StringValue(res.Namespace.ID))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *NamespaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan namespace.ResourceModel
	var state namespace.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body, _ := namespace.Converter(&plan, &state, commons.UpdateConverter)

	_, err := r.client.Client.namespace.UpdateNamespace(ctx, id, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Namespace update failed",
			fmt.Sprintf("failed to update namespace %s, error: %s", id, err.Error()),
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

func (r *NamespaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state namespace.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	_, err := r.client.Client.namespace.DeleteNamespace(ctx, id)

	if err != nil {
		errMsg := err.Error()
		resp.Diagnostics.AddError(
			"Namespace deletion failed",
			fmt.Sprintf("Failed to delete namespace %s, error: %s", id, errMsg),
		)
		return
	}
}

func (r *NamespaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

//region Private

//endregion
