package provider

import (
	"context"
	"fmt"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkStack "github.com/control-monkey/controlmonkey-sdk-go/services/stack"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/stack"
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
				MarkdownDescription: fmt.Sprintf("IaC type of the stack. Allowed values: %s.", helpers.EnumForDocs(stack.IacTypes)),
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(stack.IacTypes...),
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
				Validators: []validator.String{
					stringvalidator.NoneOf(""),
				},
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
						Validators: []validator.String{
							stringvalidator.NoneOf(""),
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
										MarkdownDescription: fmt.Sprintf("The type of the ttl. Allowed values: %s.", helpers.EnumForDocs(stack.TtlTypes)),
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf(stack.TtlTypes...),
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
	var state stack.ResourceModel
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

	res, err := r.client.Client.stack.CreateStack(ctx, &sdkStack.CreateStackInput{Stack: body})
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

func (r *StackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan stack.ResourceModel
	var state stack.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body, _ := stack.Converter(&plan, &state, commons.UpdateConverter)

	_, err := r.client.Client.stack.UpdateStack(ctx, id, body)
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
	var state stack.ResourceModel
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
