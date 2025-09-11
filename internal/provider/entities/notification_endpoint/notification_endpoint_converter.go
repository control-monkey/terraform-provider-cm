package notfication_endpoint

import (
	sdkNotification "github.com/control-monkey/controlmonkey-sdk-go/services/notification"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*sdkNotification.Endpoint, bool) {
	var retVal *sdkNotification.Endpoint

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(sdkNotification.Endpoint)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if plan.Name != state.Name {
		retVal.SetName(plan.Name.ValueStringPointer())
		hasChanges = true
	}
	if plan.Protocol != state.Protocol {
		retVal.SetProtocol(plan.Protocol.ValueStringPointer())
		hasChanges = true
	}
	if plan.Url != state.Url {
		retVal.SetUrl(plan.Url.ValueStringPointer())
		hasChanges = true
	}

	if cfg, changed := slackAppConfigConverter(plan.SlackAppConfig, state.SlackAppConfig, converterType); changed {
		retVal.SetNotificationEndpointSlackAppConfig(cfg)
		hasChanges = true
	}

	if innerProperty, hasInnerChanges := helpers.TfListStringConverter(plan.EmailAddresses, state.EmailAddresses); hasInnerChanges {
		retVal.SetEmailAddresses(innerProperty)
		hasChanges = true
	}

	return retVal, hasChanges
}

func slackAppConfigConverter(plan *SlackAppConfigModel, state *SlackAppConfigModel, converterType commons.ConverterType) (*sdkNotification.NotificationEndpointSlackAppConfig, bool) {
	var retVal *sdkNotification.NotificationEndpointSlackAppConfig

	if plan == nil {
		if state == nil {
			return nil, false
		} else {
			return nil, true
		}
	}

	retVal = new(sdkNotification.NotificationEndpointSlackAppConfig)
	hasChanges := false

	if state == nil {
		state = new(SlackAppConfigModel)
		hasChanges = true
	}

	if plan.NotificationSlackAppId != state.NotificationSlackAppId {
		retVal.SetNotificationSlackAppId(plan.NotificationSlackAppId.ValueStringPointer())
		hasChanges = true
	}
	if plan.ChannelId != state.ChannelId {
		retVal.SetChannelId(plan.ChannelId.ValueStringPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}
