package events_subscriptions

import (
	sdkNotification "github.com/control-monkey/controlmonkey-sdk-go/services/notification"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(apiEntities []*sdkNotification.EventSubscription, state *ResourceModel, scope string, scopeId *string) {
	state.Scope = helpers.StringValueOrNull(&scope)
	state.ScopeId = helpers.StringValueOrNull(scopeId)

	if apiEntities != nil {
		subscriptions := updateStateAfterReadSubscriptions(apiEntities)
		state.Subscriptions = subscriptions
	} else {
		state.Subscriptions = nil
	}
}

func updateStateAfterReadSubscriptions(permissions []*sdkNotification.EventSubscription) []*SubscriptionModel {
	var retVal []*SubscriptionModel

	if len(permissions) > 0 {
		retVal = make([]*SubscriptionModel, len(permissions))

		for i, permission := range permissions {
			u := updateStateAfterReadEventSubscription(permission)
			retVal[i] = &u
		}
	}

	return retVal
}

func updateStateAfterReadEventSubscription(permission *sdkNotification.EventSubscription) SubscriptionModel {
	var retVal SubscriptionModel

	retVal.ID = helpers.StringValueOrNull(permission.ID)
	retVal.EventType = helpers.StringValueOrNull(permission.EventType)
	retVal.NotificationEndpointId = helpers.StringValueOrNull(permission.NotificationEndpointId)

	return retVal
}
