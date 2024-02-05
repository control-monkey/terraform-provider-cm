package provider

import (
	"context"
	"fmt"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/template"
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
var _ resource.Resource = &TemplateResource{}

func NewTemplateResource() resource.Resource {
	return &TemplateResource{}
}

type TemplateResource struct {
	client *ControlMonkeyAPIClient
}

func (r *TemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

func (r *TemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys templates for ephemeral stack.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the template.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the template.",
				Required:            true,
			},
			"iac_type": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("IaC type of the template. Allowed values: %s.", helpers.EnumForDocs(cmTypes.IacTypes)),
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(cmTypes.IacTypes...),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the template.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.NoneOf(""),
				},
			},
			"vcs_info": schema.SingleNestedAttribute{
				MarkdownDescription: "The configuration of the version control to which the template is attached.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"provider_id": schema.StringAttribute{
						MarkdownDescription: "The ControlMonkey unique ID of the connected version control system.",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
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
						MarkdownDescription: "The branch that triggers the deployment of the ephemeral stack from the template. If no branch is specified, the default branch of the repository will be used.",
						Optional:            true,
					},
				},
			},
			"policy": schema.SingleNestedAttribute{
				MarkdownDescription: "The policy of the template.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"ttl_config": schema.SingleNestedAttribute{
						MarkdownDescription: "The time to live config of the template policy.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"max_ttl": schema.SingleNestedAttribute{
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
							"default_ttl": schema.SingleNestedAttribute{
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
			"skip_state_refresh_on_destroy": schema.BoolAttribute{
				MarkdownDescription: "When enabled, the state will not get refreshed before planning the destroy operation.",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *TemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *TemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get current state
	var state template.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.template.ReadTemplate(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read template %s", id), fmt.Sprintf("%s", err))
		return
	}

	template.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *TemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan template.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := template.Converter(&plan, nil, commons.CreateConverter)

	res, err := r.client.Client.template.CreateTemplate(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Template creation failed",
			fmt.Sprintf("failed to create template, error: %s", err.Error()),
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

func (r *TemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan template.ResourceModel
	var state template.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body, _ := template.Converter(&plan, &state, commons.UpdateConverter)

	_, err := r.client.Client.template.UpdateTemplate(ctx, id, body)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Template update failed",
			fmt.Sprintf("failed to update template %s, error: %s", id, err),
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

func (r *TemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state template.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	_, err := r.client.Client.template.DeleteTemplate(ctx, id)

	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Template deletion failed",
			fmt.Sprintf("Failed to delete template %s, error: %s", id, err),
		)
		return
	}
}

func (r *TemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
