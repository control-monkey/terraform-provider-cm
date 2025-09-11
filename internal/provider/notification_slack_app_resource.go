package provider

import (
	"context"
	"fmt"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	tfSlackApp "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/notification_slack_app"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &NotificationSlackAppResource{}

func NewNotificationSlackAppResource() resource.Resource { return &NotificationSlackAppResource{} }

type NotificationSlackAppResource struct{ client *ControlMonkeyAPIClient }

func (r *NotificationSlackAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_slack_app"
}

func (r *NotificationSlackAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys Slack App notifications integration. For more information: [ControlMonkey Documentation](https://docs.controlmonkey.io/administration/notifications/creating-a-slack-app)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the Slack App.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Slack App.",
				Required:            true,
				Validators:          []validator.String{cmStringValidators.NotBlank()},
			},
			"bot_auth_token": schema.StringAttribute{
				MarkdownDescription: "A sensitive bot auth token",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *NotificationSlackAppResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationSlackAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tfSlackApp.ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.notification.ListNotificationSlackApps(ctx, &id, nil)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read notification slack app", err.Error())
		return
	}
	if len(res) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	tfSlackApp.UpdateStateAfterRead(res[0], &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *NotificationSlackAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tfSlackApp.ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := tfSlackApp.Converter(&plan, nil, commons.CreateConverter)
	res, err := r.client.Client.notification.CreateNotificationSlackApp(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(resourceCreationFailedError, fmt.Sprintf("failed to create notification slack app, error: %s", err))
		return
	}

	plan.ID = types.StringValue(*res.ID)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *NotificationSlackAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tfSlackApp.ResourceModel
	var state tfSlackApp.ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	body, _ := tfSlackApp.Converter(&plan, &state, commons.UpdateConverter)
	_, err := r.client.Client.notification.UpdateNotificationSlackApp(ctx, id, body)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.Diagnostics.AddError(resourceNotFoundError, "Notification Slack App not found")
			return
		}
		resp.Diagnostics.AddError(resourceUpdateFailedError, fmt.Sprintf("failed to update notification slack app, error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *NotificationSlackAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tfSlackApp.ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := state.ID.ValueString()
	if _, err := r.client.Client.notification.DeleteNotificationSlackApp(ctx, id); err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(resourceDeletionFailedError, fmt.Sprintf("Failed to delete notification slack app, error: %s", err))
		return
	}
}

func (r *NotificationSlackAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
