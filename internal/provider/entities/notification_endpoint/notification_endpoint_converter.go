package notfication_endpoint

import (
	sdkNotification "github.com/control-monkey/controlmonkey-sdk-go/services/notification"
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

	return retVal, hasChanges
}
