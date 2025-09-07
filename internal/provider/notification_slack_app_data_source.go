package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	tfSlackAppData "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/notification_slack_app_data"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSource = &NotificationSlackAppDataSource{}

func NewNotificationSlackAppDataSource() datasource.DataSource {
	return &NotificationSlackAppDataSource{}
}

type NotificationSlackAppDataSource struct{ client *ControlMonkeyAPIClient }

func (r *NotificationSlackAppDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_slack_app"
}

func (r *NotificationSlackAppDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the slack app.",
				Optional:            true,
				Validators: []validator.String{stringvalidator.AtLeastOneOf(
					path.MatchRoot("id"), path.MatchRoot("name")),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the slack app.",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *NotificationSlackAppDataSource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationSlackAppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tfSlackAppData.ResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appId := state.ID.ValueStringPointer()
	name := state.Name.ValueStringPointer()
	res, err := r.client.Client.notification.ListNotificationSlackApps(ctx, appId, name)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read notification slack app", fmt.Sprintf("%s", err))
		return
	}
	if len(res) == 0 {
		resp.Diagnostics.AddError(resourceNotFoundError, "Notification slack app not found")
		return
	}
	if len(res) > 1 {
		resp.Diagnostics.AddError(multipleEntitiesError, fmt.Sprintf("Found multiple slack apps with name '%s'", state.Name.ValueString()))
		return
	}

	tfSlackAppData.UpdateStateAfterRead(res[0], &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
