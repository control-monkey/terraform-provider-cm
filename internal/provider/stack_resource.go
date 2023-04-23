package provider

import (
	"context"
	"fmt"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	"github.com/control-monkey/controlmonkey-sdk-go/services/stack"
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
var _ resource.Resource = &StackResource{}

func NewStackResource() resource.Resource {
	return &StackResource{}
}

type StackResource struct {
	client *ControlMonkeyAPIClient
}

type StackResourceModel struct {
	ID                 types.String             `tfsdk:"id"`
	IacType            types.String             `tfsdk:"iac_type"`
	NamespaceId        types.String             `tfsdk:"namespace_id"`
	Name               types.String             `tfsdk:"name"`
	Description        types.String             `tfsdk:"description"`
	DeploymentBehavior *deploymentBehaviorModel `tfsdk:"deployment_behavior"`
	VcsInfo            *vcsInfoModel            `tfsdk:"vcs_info"`
	RunTrigger         *runTriggerModel         `tfsdk:"run_trigger"`
	IacConfig          *IacConfigModel          `tfsdk:"iac_config"`
	Policy             *stackPolicyModel        `tfsdk:"policy"`
}

type deploymentBehaviorModel struct {
	DeployOnPush    types.Bool `tfsdk:"deploy_on_push"`
	WaitForApproval types.Bool `tfsdk:"wait_for_approval"`
}

type vcsInfoModel struct {
	ProviderId types.String `tfsdk:"provider_id"`
	RepoName   types.String `tfsdk:"repo_name"`
	Path       types.String `tfsdk:"path"`
	Branch     types.String `tfsdk:"branch"`
}

type runTriggerModel struct {
	Patterns []types.String `tfsdk:"patterns"`
}

type IacConfigModel struct {
	TerraformVersion  types.String `tfsdk:"terraform_version"`
	TerragruntVersion types.String `tfsdk:"terragrunt_version"`
}

type stackPolicyModel struct {
	TtlConfig *stackTtlConfigModel `tfsdk:"ttl_config"`
}

type stackTtlConfigModel struct {
	Ttl *stackTtlDefinitionModel `tfsdk:"ttl"`
}

type stackTtlDefinitionModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.Int64  `tfsdk:"value"`
}

var iacTypes = []string{"terraform", "terragrunt"}
var stackTtlTypes = []string{"hours", "days"}

func (r *StackResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stack"
}

func (r *StackResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys stacks.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the stack.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"iac_type": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("IaC type of the stack. Allowed values: %s.", helpers.EnumForDocs(iacTypes)),
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(iacTypes...),
				},
			},
			"namespace_id": schema.StringAttribute{
				MarkdownDescription: "The id of the namespace that contains the stack.",
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
			},
			"deployment_behavior": schema.SingleNestedAttribute{
				MarkdownDescription: "The deployment behavior configuration.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"deploy_on_push": schema.BoolAttribute{
						MarkdownDescription: "Whether to start a deployment on a push event or not.",
						Required:            true,
					},
					"wait_for_approval": schema.BoolAttribute{
						MarkdownDescription: "Whether to wait to a manual approval before deployment or not.",
						Required:            true,
					},
				},
			},
			"vcs_info": schema.SingleNestedAttribute{
				MarkdownDescription: "The configuration of the version control to which the stack is attached.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"provider_id": schema.StringAttribute{
						MarkdownDescription: "The Control Monkey unique ID of the connected version control system.",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"repo_name": schema.StringAttribute{
						MarkdownDescription: "The name of the version control repository.",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"path": schema.StringAttribute{
						MarkdownDescription: "The path to a chosen directory from the root.",
						Optional:            true,
					},
					"branch": schema.StringAttribute{
						MarkdownDescription: "The branch that should trigger plan/deployment for the stack. When no branch is given, the default branch of the repository is chosen.",
						Optional:            true,
					},
				},
			},
			"run_trigger": schema.SingleNestedAttribute{
				MarkdownDescription: "Glob patterns to specify additional paths that should trigger a stack run.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"patterns": schema.ListAttribute{
						MarkdownDescription: "Patterns that trigger a stack run.",
						ElementType:         types.StringType,
						Optional:            true,
					},
				},
			},
			"iac_config": schema.SingleNestedAttribute{
				MarkdownDescription: "IaC configuration of the stack.",
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
				},
			},
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
										MarkdownDescription: fmt.Sprintf("The type of the ttl. Allowed values: %s.", helpers.EnumForDocs(stackTtlTypes)),
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf(stackTtlTypes...),
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

// Read refreshes the Terraform state with the latest data.
func (r *StackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get current state
	var state StackResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.stack.ReadStack(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read stack %s", id), fmt.Sprintf("%s", err))
		return
	}

	stack := res.Stack

	state.IacType = helpers.StringValueOrNull(stack.IacType)
	state.NamespaceId = helpers.StringValueOrNull(stack.NamespaceId)
	state.Name = helpers.StringValueOrNull(stack.Name)
	state.Description = helpers.StringValueOrNull(stack.Description)

	data := stack.Data
	var dp deploymentBehaviorModel
	dp.DeployOnPush = helpers.BoolValueOrNull(data.DeploymentBehavior.DeployOnPush)
	dp.WaitForApproval = helpers.BoolValueOrNull(data.DeploymentBehavior.WaitForApproval)
	state.DeploymentBehavior = &dp

	var vcs vcsInfoModel
	vcs.ProviderId = helpers.StringValueOrNull(data.VcsInfo.ProviderId)
	vcs.RepoName = helpers.StringValueOrNull(data.VcsInfo.RepoName)
	vcs.Path = helpers.StringValueOrNull(data.VcsInfo.Path)
	vcs.Branch = helpers.StringValueOrNull(data.VcsInfo.Branch)
	state.VcsInfo = &vcs

	var rt runTriggerModel
	if data.RunTrigger != nil {
		if data.RunTrigger.Patterns != nil {
			rt.Patterns = helpers.StringSlice(data.RunTrigger.Patterns)
		} else {
			rt.Patterns = nil
		}
		state.RunTrigger = &rt
	} else {
		state.RunTrigger = nil
	}

	var ic IacConfigModel
	if data.IacConfig != nil {
		ic.TerraformVersion = helpers.StringValueOrNull(data.IacConfig.TerraformVersion)
		ic.TerragruntVersion = helpers.StringValueOrNull(data.IacConfig.TerragruntVersion)
		state.IacConfig = &ic
	} else {
		state.IacConfig = nil
	}

	var policy stackPolicyModel
	var ttlc stackTtlConfigModel
	var ttl stackTtlDefinitionModel
	if data.Policy != nil {
		if data.Policy.TtlConfig != nil {
			ttl.Type = helpers.StringValueOrNull(data.Policy.TtlConfig.Ttl.Type)
			ttl.Value = helpers.Int64ValueOrNull(data.Policy.TtlConfig.Ttl.Value)
			ttlc.Ttl = &ttl
			policy.TtlConfig = &ttlc
		} else {
			policy.TtlConfig = nil
		}
		state.Policy = &policy
	} else {
		state.Policy = nil
	}
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
	var plan StackResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := createCmStackFromPlanConverter(plan)

	res, err := r.client.Client.stack.CreateStack(ctx, &stack.CreateStackInput{Stack: &body})
	if err != nil {
		resp.Diagnostics.AddError(
			"Stack creation failed",
			fmt.Sprintf("failed to create stack, error: %s", err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(controlmonkey.StringValue(res.Stack.ID))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func createCmStackFromPlanConverter(plan StackResourceModel) stack.Stack {
	var body stack.Stack

	body.SetIacType(plan.IacType.ValueStringPointer())
	body.SetNamespaceId(plan.NamespaceId.ValueStringPointer())
	body.SetName(plan.Name.ValueStringPointer())
	body.SetDescription(plan.Description.ValueStringPointer())

	var data stack.Data

	var db stack.DeploymentBehavior
	db.SetDeployOnPush(plan.DeploymentBehavior.DeployOnPush.ValueBoolPointer())
	db.SetWaitForApproval(plan.DeploymentBehavior.WaitForApproval.ValueBoolPointer())
	data.SetDeploymentBehavior(&db)

	var vi stack.VcsInfo
	vi.SetProviderId(plan.VcsInfo.ProviderId.ValueStringPointer())
	vi.SetRepoName(plan.VcsInfo.RepoName.ValueStringPointer())
	vi.SetPath(plan.VcsInfo.Path.ValueStringPointer())
	vi.SetBranch(plan.VcsInfo.Branch.ValueStringPointer())
	data.SetVcsInfo(&vi)

	var rt stack.RunTrigger
	runTrigger := plan.RunTrigger

	if runTrigger != nil {
		if runTrigger.Patterns != nil {
			var patterns []*string
			for _, pattern := range plan.RunTrigger.Patterns {
				patterns = append(patterns, pattern.ValueStringPointer())
			}
			rt.SetPatterns(patterns)
		} else {
			rt.SetPatterns(nil)
		}
		data.SetRunTrigger(&rt)
	} else {
		data.SetRunTrigger(nil)
	}

	var ic stack.IacConfig
	if plan.IacConfig != nil {
		ic.SetTerraformVersion(plan.IacConfig.TerraformVersion.ValueStringPointer())
		ic.SetTerragruntVersion(plan.IacConfig.TerragruntVersion.ValueStringPointer())
		data.SetIacConfig(&ic)
	} else {
		data.SetIacConfig(nil)
	}

	var policy stack.Policy
	var ttlConfig stack.TtlConfig
	var ttl stack.TtlDefinition
	if plan.Policy != nil {
		if plan.Policy.TtlConfig != nil {
			ttl.SetType(plan.Policy.TtlConfig.Ttl.Type.ValueStringPointer())
			ttl.SetValue(controlmonkey.Int(int(plan.Policy.TtlConfig.Ttl.Value.ValueInt64())))
			ttlConfig.SetTtl(&ttl)
			policy.SetTtlConfig(&ttlConfig)
		} else {
			policy.SetTtlConfig(nil)
		}
		data.SetPolicy(&policy)
	} else {
		data.SetPolicy(nil)
	}

	body.SetData(&data)
	return body
}

// do not copy it - would be better to have both create and update on the same function with a parameter to distinguish between them
func updateCmStackFromPlanConverter(plan StackResourceModel) stack.Stack {
	var body stack.Stack

	body.SetName(plan.Name.ValueStringPointer())
	body.SetDescription(plan.Description.ValueStringPointer())

	var data stack.Data

	var db stack.DeploymentBehavior
	db.SetDeployOnPush(plan.DeploymentBehavior.DeployOnPush.ValueBoolPointer())
	db.SetWaitForApproval(plan.DeploymentBehavior.WaitForApproval.ValueBoolPointer())
	data.SetDeploymentBehavior(&db)

	var vi stack.VcsInfo
	vi.SetProviderId(plan.VcsInfo.ProviderId.ValueStringPointer())
	vi.SetRepoName(plan.VcsInfo.RepoName.ValueStringPointer())
	vi.SetPath(plan.VcsInfo.Path.ValueStringPointer())
	vi.SetBranch(plan.VcsInfo.Branch.ValueStringPointer())
	data.SetVcsInfo(&vi)

	var rt stack.RunTrigger
	if plan.RunTrigger != nil {
		if plan.RunTrigger.Patterns != nil {
			var patterns []*string
			for _, pattern := range plan.RunTrigger.Patterns {
				patterns = append(patterns, pattern.ValueStringPointer())
			}
			rt.SetPatterns(patterns)
		} else {
			rt.SetPatterns(nil)
		}
		data.SetRunTrigger(&rt)
	} else {
		data.SetRunTrigger(nil)
	}

	var ic stack.IacConfig
	if plan.IacConfig != nil {
		ic.SetTerraformVersion(plan.IacConfig.TerraformVersion.ValueStringPointer())
		ic.SetTerragruntVersion(plan.IacConfig.TerragruntVersion.ValueStringPointer())
		data.SetIacConfig(&ic)
	} else {
		data.SetIacConfig(nil)
	}

	var policy stack.Policy
	var ttlConfig stack.TtlConfig
	var ttl stack.TtlDefinition
	if plan.Policy != nil {
		if plan.Policy.TtlConfig != nil {
			ttl.SetType(plan.Policy.TtlConfig.Ttl.Type.ValueStringPointer())
			ttl.SetValue(controlmonkey.Int(int(plan.Policy.TtlConfig.Ttl.Value.ValueInt64())))
			ttlConfig.SetTtl(&ttl)
			policy.SetTtlConfig(&ttlConfig)
		} else {
			policy.SetTtlConfig(nil)
		}
		data.SetPolicy(&policy)
	} else {
		data.SetPolicy(nil)
	}

	body.SetData(&data)
	return body
}

func (r *StackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan StackResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body := updateCmStackFromPlanConverter(plan)

	_, err := r.client.Client.stack.UpdateStack(ctx, id, &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Stack update failed",
			fmt.Sprintf("failed to update stack %s, error: %s", id, err.Error()),
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
	var state StackResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	_, err := r.client.Client.stack.DeleteStack(ctx, id)

	if err != nil {
		errMsg := err.Error()
		resp.Diagnostics.AddError(
			"Stack deletion failed",
			fmt.Sprintf("Failed to delete stack %s, error: %s", id, errMsg),
		)
		return
	}
}

func (r *StackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
