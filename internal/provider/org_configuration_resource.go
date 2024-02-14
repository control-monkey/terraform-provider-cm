package provider

import (
	"context"
	"fmt"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	tfOrgConfiguration "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/org_configuration"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OrgConfigurationResource{}

func NewOrgConfigurationResource() resource.Resource {
	return &OrgConfigurationResource{}
}

type OrgConfigurationResource struct {
	client *ControlMonkeyAPIClient
}

func (r *OrgConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_configuration"
}

func (r *OrgConfigurationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys org configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of this resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"iac_config": schema.SingleNestedAttribute{
				MarkdownDescription: "IaC configuration that defines default versions. If not explicitly overridden, these defaults apply to all namespaces/stacks.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"terraform_version": schema.StringAttribute{
						MarkdownDescription: "the Terraform version that will be used for terraform operations.",
						Optional:            true,
						Validators: []validator.String{
							cmStringValidators.NotBlank(),
						},
					},
					"terragrunt_version": schema.StringAttribute{
						MarkdownDescription: "the Terragrunt version that will be used for terragrunt operations.",
						Optional:            true,
						Validators: []validator.String{
							cmStringValidators.NotBlank(),
						},
					},
					"opentofu_version": schema.StringAttribute{
						MarkdownDescription: "the OpenTofu version that will be used for tofu operations.",
						Optional:            true,
						Validators: []validator.String{
							cmStringValidators.NotBlank(),
						},
					},
				},
			},
			"runner_config": schema.SingleNestedAttribute{
				MarkdownDescription: "Configure the runner settings to specify whether ControlMonkey manages the runner or it is self-hosted.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"mode": schema.StringAttribute{
						MarkdownDescription: fmt.Sprintf("The runner mode. Allowed values: %s.", helpers.EnumForDocs(cmTypes.RunnerConfigModeTypes)),
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(cmTypes.RunnerConfigModeTypes...),
						},
					},
					"groups": schema.ListAttribute{
						MarkdownDescription: fmt.Sprintf("In case that `mode` is `%s`, groups must contain at least one runners group. If `mode` is `%s`, this field must not be configures.", cmTypes.SelfHosted, cmTypes.Managed),
						ElementType:         types.StringType,
						Optional:            true,
						Validators:          commons.ValidateUniqueListWithNoBlankValues(),
					},
					"is_overridable": schema.BoolAttribute{
						MarkdownDescription: "By setting this option, you allow this configuration to be overridden in specific namespaces/stacks.",
						Optional:            true,
					},
				},
			},
			"s3_state_files_locations": schema.ListNestedAttribute{
				MarkdownDescription: "The S3 buckets of your current terraform state files. This will be used by ControlMonkey to scan for existing managed resources.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"bucket_name": schema.StringAttribute{
							MarkdownDescription: "The name of the bucket.",
							Required:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
						"bucket_region": schema.StringAttribute{
							MarkdownDescription: "The region of the bucket.",
							Required:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
						"aws_account_id": schema.StringAttribute{
							MarkdownDescription: "The AWS account ID in which the bucket is situated.",
							Required:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
					},
				},
			},
			"suppressed_resources": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"managed_by_tags": schema.ListNestedAttribute{
						MarkdownDescription: "List of tags by which any AWS resource with one of the configured tags will be considered as managed. The tag key/value definition is case sensitive.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									MarkdownDescription: "The key of the tag.",
									Required:            true,
									Validators: []validator.String{
										cmStringValidators.NotBlank(),
									},
								},
								"value": schema.StringAttribute{
									MarkdownDescription: "The value of the tag.",
									Optional:            true,
									Validators: []validator.String{
										cmStringValidators.NotBlank(),
									},
								},
							},
						},
					},
				},
			},
			"report_configurations": schema.ListNestedAttribute{
				MarkdownDescription: "The S3 buckets of your current terraform state files. This will be used by ControlMonkey to scan for existing managed resources.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Indicates whether the report distribution is enabled or disabled.",
							Required:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: fmt.Sprintf("The type of the report. Supported types: %s.", helpers.EnumForDocs(cmTypes.ReportTypes)),
							Required:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
						"recipients": schema.SingleNestedAttribute{
							MarkdownDescription: "Specifies who will receive the report.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"all_admins": schema.BoolAttribute{
									MarkdownDescription: "If enabled, the report will be sent to every administrator within your organization.",
									Optional:            true,
								},
								"email_addresses": schema.ListAttribute{
									MarkdownDescription: "List of email addresses to which the report will be sent.",
									ElementType:         types.StringType,
									Optional:            true,
									Validators:          commons.ValidateUniqueListWithNoBlankValues(),
								},
								"email_addresses_to_exclude": schema.ListAttribute{
									MarkdownDescription: "List of email addresses to which the report will not be sent.",
									ElementType:         types.StringType,
									Optional:            true,
									Validators:          commons.ValidateUniqueListWithNoBlankValues(),
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
func (r *OrgConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *OrgConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get current state
	var state tfOrgConfiguration.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Client.organization.ReadOrgConfiguration(ctx)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read org configuration"), err.Error())
		return
	}

	tfOrgConfiguration.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *OrgConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//only a single org config must exist. So, before creating a new one, we check if one is already exists
	diags := r.checkIfExistsBeforeCreate(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Retrieve values from plan
	var plan tfOrgConfiguration.ResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := tfOrgConfiguration.Converter(&plan, nil, commons.CreateConverter)

	if _, err := r.client.Client.organization.UpsertOrgConfiguration(ctx, body); err != nil {
		resp.Diagnostics.AddError(
			resourceCreationFailedError,
			fmt.Sprintf("Failed to create org configuration, error: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(tfOrgConfiguration.ImportID)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *OrgConfigurationResource) checkIfExistsBeforeCreate(ctx context.Context) diag.Diagnostics {
	retVal := diag.Diagnostics{}

	res, err := r.client.Client.organization.ReadOrgConfiguration(ctx)

	if err != nil {
		retVal.AddError(resourceCreationFailedError, fmt.Sprintf("Failed to create org configuration. Error: %s", err))
	} else if helpers.IsAllNilFields(res) == false {
		retVal.AddError("Org Configuration already exists, there is only one configuration allowed per organization",
			fmt.Sprintf("Import operation is required to manage this resourcce. Use import command e.g 'terraform import cm_org_configuration.<resource_name> %s'", tfOrgConfiguration.ImportID),
		)
	}

	return retVal
}

func (r *OrgConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan tfOrgConfiguration.ResourceModel
	var state tfOrgConfiguration.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := tfOrgConfiguration.Converter(&plan, &state, commons.UpdateConverter)

	if _, err := r.client.Client.organization.UpsertOrgConfiguration(ctx, body); err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.Diagnostics.AddError(resourceNotFoundError, fmt.Sprintf("Org Configuration not found"))
			return
		}
		resp.Diagnostics.AddError(
			resourceUpdateFailedError,
			fmt.Sprintf("failed to update org configuration, error: %s", err),
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

func (r *OrgConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state tfOrgConfiguration.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if _, err := r.client.Client.organization.DeleteOrgConfiguration(ctx); err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			resourceDeletionFailedError,
			fmt.Sprintf("Failed to delete org configuration, error: %s", err),
		)
		return
	}
}

func (r *OrgConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != tfOrgConfiguration.ImportID {
		resp.Diagnostics.AddError(validationError, fmt.Sprintf("ID must be '%s'", tfOrgConfiguration.ImportID))
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
