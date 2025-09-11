package notfication_endpoint

import (
	sdkNotification "github.com/control-monkey/controlmonkey-sdk-go/services/notification"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res *sdkNotification.Endpoint, state *ResourceModel) {
	state.Name = helpers.StringValueOrNull(res.Name)
	state.Protocol = helpers.StringValueOrNull(res.Protocol)
	state.Url = helpers.StringValueOrNull(res.Url)

	if res.NotificationEndpointSlackAppConfig != nil {
		cfg := updateStateAfterReadSlackAppConfig(res.NotificationEndpointSlackAppConfig)
		state.SlackAppConfig = &cfg
	} else {
		state.SlackAppConfig = nil
	}

	state.EmailAddresses = helpers.StringPointerSliceToTfList(res.EmailAddresses)
}

func updateStateAfterReadSlackAppConfig(cfg *sdkNotification.NotificationEndpointSlackAppConfig) SlackAppConfigModel {
	var retVal SlackAppConfigModel

	retVal.NotificationSlackAppId = helpers.StringValueOrNull(cfg.NotificationSlackAppId)
	retVal.ChannelId = helpers.StringValueOrNull(cfg.ChannelId)

	return retVal
}
