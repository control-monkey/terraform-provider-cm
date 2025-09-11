package notfication_endpoint

import "github.com/hashicorp/terraform-plugin-framework/types"

type ResourceModel struct {
	ID             types.String         `tfsdk:"id"`
	Name           types.String         `tfsdk:"name"`
	Protocol       types.String         `tfsdk:"protocol"`
	Url            types.String         `tfsdk:"url"`
	SlackAppConfig *SlackAppConfigModel `tfsdk:"slack_app_config"`
	EmailAddresses types.List           `tfsdk:"email_addresses"`
}

type SlackAppConfigModel struct {
	NotificationSlackAppId types.String `tfsdk:"notification_slack_app_id"`
	ChannelId              types.String `tfsdk:"channel_id"`
}
