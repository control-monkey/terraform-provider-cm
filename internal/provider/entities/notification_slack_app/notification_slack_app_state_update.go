package notification_slack_app

import (
	api "github.com/control-monkey/controlmonkey-sdk-go/services/notification"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res *api.NotificationSlackApp, state *ResourceModel) {
	state.ID = helpers.StringValueOrNull(res.ID)
	state.Name = helpers.StringValueOrNull(res.Name)
}
