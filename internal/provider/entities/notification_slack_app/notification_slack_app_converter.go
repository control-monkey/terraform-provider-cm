package notification_slack_app

import (
	api "github.com/control-monkey/controlmonkey-sdk-go/services/notification"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*api.NotificationSlackApp, bool) {
	var retVal *api.NotificationSlackApp

	if plan == nil {
		if state == nil {
			return nil, false
		} else {
			return nil, true
		}
	}

	retVal = new(api.NotificationSlackApp)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel)
		hasChanges = true
	}

	if plan.Name != state.Name {
		retVal.SetName(plan.Name.ValueStringPointer())
		hasChanges = true
	}

	if plan.BotAuthToken != state.BotAuthToken {
		retVal.SetBotAuthToken(plan.BotAuthToken.ValueStringPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}
