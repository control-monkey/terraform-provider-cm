package notification_slack_app

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	BotAuthToken types.String `tfsdk:"bot_auth_token"`
}
