package provider

import (
	"context"
	"fmt"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	tfDisasterRecoveryConfiguration "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/disaster_recovery_configuration"
	cm_stringvalidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
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
var _ resource.Resource = &DisasterRecoveryConfigurationResource{}

func NewDisasterRecoveryConfigurationResource() resource.Resource {
	return &DisasterRecoveryConfigurationResource{}
}

type DisasterRecoveryConfigurationResource struct {
	client *ControlMonkeyAPIClient
}

func (r *DisasterRecoveryConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_disaster_recovery_configuration"
}

func (r *DisasterRecoveryConfigurationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys disaster recovery configurations. For more information: [ControlMonkey Documentation](https://docs.controlmonkey.io/main-concepts/disaster-recovery)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the disaster recovery configuration.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "Specifies the cloud provider type, such as `aws`.",
				Required:            true,
				Validators:          []validator.String{cm_stringvalidators.NotBlank()},
			},
			"cloud_account_id": schema.StringAttribute{
				MarkdownDescription: "The identifier of the cloud account, such as an AWS Account ID for AWS or a Subscription ID for Azure.",
				Required:            true,
				Validators:          []validator.String{cm_stringvalidators.NotBlank()},
			},
			"backup_strategy": schema.SingleNestedAttribute{
				MarkdownDescription: "The configuration specifying where the backup will be stored and how to break down your resources.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"include_managed_resources": schema.BoolAttribute{
						MarkdownDescription: "Indicates whether managed resources should be included in the backup.",
						Required:            true,
					},
					"mode": schema.StringAttribute{
						MarkdownDescription: fmt.Sprintf("Specify the backup strategy mode whether you want the default ControlMonkey behaviour or your own custom strategy. Allowed values: %s.", helpers.EnumForDocs(cmTypes.DisasterRecoveryBackupModeTypes)),
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(cmTypes.DisasterRecoveryBackupModeTypes...),
							cm_stringvalidators.NotBlank(),
						},
					},
					"vcs_info": schema.SingleNestedAttribute{
						MarkdownDescription: "Configuration details for the version control system where the disaster recovery snapshot will be stored.",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"provider_id": schema.StringAttribute{
								MarkdownDescription: "The ControlMonkey unique ID of the connected version control system.",
								Required:            true,
								Validators:          []validator.String{cm_stringvalidators.NotBlank()},
							},
							"repo_name": schema.StringAttribute{
								MarkdownDescription: "The name of the version control repository.",
								Required:            true,
								Validators:          []validator.String{cm_stringvalidators.NotBlank()},
							},
							"branch": schema.StringAttribute{
								MarkdownDescription: "The target branch the disaster recovery snapshot will be pushed to.",
								Required:            true,
								Validators:          []validator.String{cm_stringvalidators.NotBlank()},
							},
						},
					},
					"groups_json": schema.StringAttribute{
						MarkdownDescription: fmt.Sprintf("JSON format of your custom strategy. Describe how to group the resources we backup into your VCS. This field is required only when `mode` is set to `manual`. This filed is not allowed when `mode` is set to `default`.\nFor more information: [ControlMonkey Documentation](https://docs.controlmonkey.io/main-concepts/disaster-recovery/infrastructure-daily-backup#how-to-configure)"),
						Optional:            true,
						CustomType:          jsontypes.NormalizedType{},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *DisasterRecoveryConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *DisasterRecoveryConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state tfDisasterRecoveryConfiguration.ResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.disasterRecovery.ReadDisasterRecoveryConfiguration(ctx, id)

	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(resourceNotFoundError, fmt.Sprintf("Disaster Recovery Configuration '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read disaster recovery configuration '%s'", id), err.Error())
		return
	}

	tfDisasterRecoveryConfiguration.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *DisasterRecoveryConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan tfDisasterRecoveryConfiguration.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := tfDisasterRecoveryConfiguration.Converter(&plan, nil, commons.CreateConverter)

	res, err := r.client.Client.disasterRecovery.CreateDisasterRecoveryConfiguration(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceCreationFailedError,
			fmt.Sprintf("failed to create disaster recovery configuration, error: %s", err.Error()),
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

func (r *DisasterRecoveryConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan tfDisasterRecoveryConfiguration.ResourceModel
	var state tfDisasterRecoveryConfiguration.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body, _ := tfDisasterRecoveryConfiguration.Converter(&plan, &state, commons.UpdateConverter)

	_, err := r.client.Client.disasterRecovery.UpdateDisasterRecoveryConfiguration(ctx, id, body)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.Diagnostics.AddError(resourceNotFoundError, fmt.Sprintf("Disaster Recovery Configuration '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(
			resourceUpdateFailedError,
			fmt.Sprintf("failed to update disaster recovery configuration %s, error: %s", id, err),
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

func (r *DisasterRecoveryConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state tfDisasterRecoveryConfiguration.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	_, err := r.client.Client.disasterRecovery.DeleteDisasterRecoveryConfiguration(ctx, id)

	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			resourceDeletionFailedError,
			fmt.Sprintf("Failed to delete disaster recovery configuration %s, error: %s", id, err),
		)
		return
	}
}

func (r *DisasterRecoveryConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
