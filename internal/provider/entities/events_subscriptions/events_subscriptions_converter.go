package events_subscriptions

import (
	sdkNotification "github.com/control-monkey/controlmonkey-sdk-go/services/notification"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	"github.com/hashicorp/go-set/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MergedEntities struct {
	EntitiesToCreate []*sdkNotification.EventSubscription
	EntitiesToUpdate []*sdkNotification.EventSubscription
	EntitiesToDelete []*sdkNotification.EventSubscription
}

func Merge(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) *MergedEntities {
	retVal := new(MergedEntities)

	if plan == nil {
		plan = new(ResourceModel) // delete merger
	}

	if state == nil {
		state = new(ResourceModel) // create merger
	}

	var scope types.String
	var scopeId types.String
	if plan.Scope.IsNull() == false {
		scope = plan.Scope
		scopeId = plan.ScopeId
	} else {
		scope = state.Scope
		scopeId = plan.ScopeId
	}

	mergeResult := interfaces.MergeEntities(plan.Subscriptions, state.Subscriptions)
	retVal.EntitiesToCreate = convertEntities(mergeResult.EntitiesToCreate, scope, scopeId, interfaces.CreateOperation)
	retVal.EntitiesToDelete = convertEntities(mergeResult.EntitiesToDelete, scope, scopeId, interfaces.DeleteOperation)

	return retVal
}

func convertEntities(entities set.Collection[*SubscriptionModel], scope types.String, scopeId types.String, operation interfaces.OperationType) []*sdkNotification.EventSubscription {
	retVal := make([]*sdkNotification.EventSubscription, entities.Size())

	for i, e := range entities.Slice() {
		apiEntity := new(sdkNotification.EventSubscription)

		if operation == interfaces.CreateOperation {
			apiEntity.SetNotificationEndpointId(e.NotificationEndpointId.ValueStringPointer())
			apiEntity.SetScope(scope.ValueStringPointer())
			apiEntity.SetScopeId(scopeId.ValueStringPointer())
			apiEntity.SetEventType(e.EventType.ValueStringPointer())
		} else { //operation == interfaces.DeleteOperation
			apiEntity.SetID(e.ID.ValueStringPointer())
		}

		retVal[i] = apiEntity
	}

	return retVal
}
