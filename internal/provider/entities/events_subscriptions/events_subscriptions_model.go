package events_subscriptions

import (
	"fmt"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type ResourceModel struct {
	ID            types.String         `tfsdk:"id"`
	Scope         types.String         `tfsdk:"scope"`
	ScopeId       types.String         `tfsdk:"scope_id"`
	Subscriptions []*SubscriptionModel `tfsdk:"subscriptions"`
}

type SubscriptionModel struct { //When new field is added consider Hash() function
	ID                     types.String `tfsdk:"id"`
	EventType              types.String `tfsdk:"event_type"`
	NotificationEndpointId types.String `tfsdk:"notification_endpoint_id"`
}

func (e *SubscriptionModel) Hash() string {
	retVal := ""

	retVal += fmt.Sprintf("EventType:%s:", e.EventType.ValueString())
	retVal += fmt.Sprintf("NotificationEndpointId:%s:", e.NotificationEndpointId.ValueString())

	return retVal
}

func (e *SubscriptionModel) GetBlockIdentifier() string {
	retVal := ""

	if helpers.IsKnown(e.EventType) && helpers.IsKnown(e.NotificationEndpointId) {
		retVal += fmt.Sprintf("EventType:%s:", e.EventType.ValueString())
		retVal += fmt.Sprintf("NotificationEndpointId:%s:", e.NotificationEndpointId.ValueString())
	}

	return retVal
}

func CleanIdentifier(s string) string {
	split := strings.Split(s, ":")
	return fmt.Sprintf("%s %s", split[1], split[3])

}
