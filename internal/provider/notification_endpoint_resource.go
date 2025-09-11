package provider

import (
	"context"
	"fmt"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	tfNotificationEndpoint "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/notification_endpoint"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
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
var _ resource.Resource = &NotificationEndpointResource{}

func NewNotificationEndpointResource() resource.Resource {
	return &NotificationEndpointResource{}
}

type NotificationEndpointResource struct {
	client *ControlMonkeyAPIClient
}

func (r *NotificationEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_endpoint"
}

func (r *NotificationEndpointResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys notification endpoints. For more information: [ControlMonkey Documentation](https://docs.controlmonkey.io/administration/notifications)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the endpoint.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the endpoint.",
				Required:            true,
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The approach to publish notifications. Allowed values: %s.", helpers.EnumForDocs(cmTypes.EventSubscriptionProtocolTypes)),
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(cmTypes.EventSubscriptionProtocolTypes...),
				},
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The webhook url to which the notification will be sent. Required when `protocol` is one of [**slack**, **teams**]. Conflicts with `email_addresses` and `slack_app_config`.",
				Optional:            true,
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
					stringvalidator.ExactlyOneOf(path.MatchRoot("url"), path.MatchRoot("email_addresses"), path.MatchRoot("slack_app_config")),
				},
			},
			"email_addresses": schema.ListAttribute{
				MarkdownDescription: "List of email addresses to notify. Required when `protocol` is **email**. Conflicts with `url` and `slack_app_config`.",
				Optional:            true,
				ElementType:         types.StringType,
				Validators:          commons.ValidateUniqueNotEmptyListWithNoBlankValues(),
			},
			"slack_app_config": schema.SingleNestedAttribute{
				MarkdownDescription: "Slack App configuration. Required when `protocol` is **slackApp**. Conflicts with `email_addresses` and `url`.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"notification_slack_app_id": schema.StringAttribute{
						MarkdownDescription: "The Slack App ID.",
						Required:            true,
						Validators:          []validator.String{cmStringValidators.NotBlank()},
					},
					"channel_id": schema.StringAttribute{
						MarkdownDescription: "The Slack channel ID.",
						Required:            true,
						Validators:          []validator.String{cmStringValidators.NotBlank()},
					},
				},
			},
		},
	}
}

// ValidateConfig enforces conditional requirements based on selected protocol
func (r *NotificationEndpointResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data tfNotificationEndpoint.ResourceModel
	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		return
	}

	if helpers.IsKnown(data.Protocol) == false {
		return
	}

	// Exactly one of url, email_addresses, slack_app_config must be configured
	presentCount := 0
	if helpers.IsKnown(data.Url) {
		presentCount++
	}
	if data.SlackAppConfig != nil {
		presentCount++
	}
	if helpers.IsKnown(data.EmailAddresses) {
		presentCount++
	}

	if presentCount == 0 {
		return
	}

	protocol := data.Protocol.ValueString()
	switch protocol {
	case cmTypes.SlackProtocol, cmTypes.TeamsProtocol:
		if helpers.IsKnown(data.Url) == false {
			resp.Diagnostics.AddError(validationError, fmt.Sprintf("'url' is required when protocol is '%s'", protocol))
		}
	case cmTypes.SlackAppProtocol:
		if data.SlackAppConfig == nil {
			resp.Diagnostics.AddError(validationError, "'slack_app_config' is required when protocol is 'slackApp'")
		}
	case cmTypes.EmailProtocol:
		if helpers.IsKnown(data.EmailAddresses) == false {
			resp.Diagnostics.AddError(validationError, "'email_addresses' is required when protocol is 'email'")
		}
	}
}

// Configure adds the provider configured client to the data source.
func (r *NotificationEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *NotificationEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get current state
	var state tfNotificationEndpoint.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.notification.ReadNotificationEndpoint(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read notification endpoint '%s'", id), err.Error())
		return
	}

	tfNotificationEndpoint.UpdateStateAfterRead(res, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *NotificationEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan tfNotificationEndpoint.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := tfNotificationEndpoint.Converter(&plan, nil, commons.CreateConverter)

	res, err := r.client.Client.notification.CreateNotificationEndpoint(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceCreationFailedError,
			fmt.Sprintf("failed to create notification endpoint, error: %s", err.Error()),
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

func (r *NotificationEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan tfNotificationEndpoint.ResourceModel
	var state tfNotificationEndpoint.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body, _ := tfNotificationEndpoint.Converter(&plan, &state, commons.UpdateConverter)

	_, err := r.client.Client.notification.UpdateNotificationEndpoint(ctx, id, body)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.Diagnostics.AddError(resourceNotFoundError, fmt.Sprintf("Notification endpoint '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(
			resourceUpdateFailedError, fmt.Sprintf("failed to update notification endpoint %s, error: %s", id, err.Error()),
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

func (r *NotificationEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state tfNotificationEndpoint.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	_, err := r.client.Client.notification.DeleteNotificationEndpoint(ctx, id)

	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			resourceDeletionFailedError,
			fmt.Sprintf("Failed to delete notification endpoint %s, error: %s", id, err),
		)
		return
	}
}

func (r *NotificationEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
