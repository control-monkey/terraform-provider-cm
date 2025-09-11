package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	tfStackDependency "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/stack_dependency"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &StackDependencyResource{}

func NewStackDependencyResource() resource.Resource { return &StackDependencyResource{} }

type StackDependencyResource struct{ client *ControlMonkeyAPIClient }

func (r *StackDependencyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stack_dependency"
}

func (r *StackDependencyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys a stack dependency. For more information: [ControlMonkey Documentation](https://docs.controlmonkey.io/main-concepts/stack/stack-dependencies)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the dependency.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"stack_id": schema.StringAttribute{
				MarkdownDescription: "The stack that depends on another stack.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"depends_on_stack_id": schema.StringAttribute{
				MarkdownDescription: "The stack to depend on.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"trigger_option": schema.StringAttribute{
				MarkdownDescription: "Dependency trigger option. When set, `references` is required. Find supported types [here](https://docs.controlmonkey.io/controlmonkey-api/api-enumerations#stack-dependency-trigger-option-types)",
				Optional:            true,
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("references"),
					}...),
				},
			},
			"references": schema.ListNestedAttribute{
				MarkdownDescription: "List of references wiring outputs to inputs.When set, `trigger_option` is required",
				Optional:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("trigger_option"),
					}...),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"output_of_stack_to_depend_on": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
						"input_for_stack": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
						"include_sensitive_output": schema.BoolAttribute{
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *StackDependencyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StackDependencyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tfStackDependency.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.stack.ReadDependency(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read dependency '%s'", id), err.Error())
		return
	}
	tfStackDependency.UpdateStateAfterRead(res, &state)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *StackDependencyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tfStackDependency.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := tfStackDependency.Converter(&plan, nil, commons.CreateConverter)
	res, err := r.client.Client.stack.CreateDependency(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(resourceCreationFailedError, fmt.Sprintf("failed to create stack dependency, error: %s", err.Error()))
		return
	}

	plan.ID = types.StringValue(controlmonkey.StringValue(res.ID))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *StackDependencyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tfStackDependency.ResourceModel
	var state tfStackDependency.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body, _ := tfStackDependency.Converter(&plan, &state, commons.UpdateConverter)
	_, err := r.client.Client.stack.UpdateDependency(ctx, id, body)
	if err != nil {
		resp.Diagnostics.AddError(resourceUpdateFailedError, fmt.Sprintf("failed to update stack dependency %s, error: %s", id, err.Error()))
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *StackDependencyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tfStackDependency.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := state.ID.ValueString()
	_, err := r.client.Client.stack.DeleteDependency(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(resourceDeletionFailedError, fmt.Sprintf("Failed to delete dependency %s, error: %s", id, err))
		return
	}
}

func (r *StackDependencyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
