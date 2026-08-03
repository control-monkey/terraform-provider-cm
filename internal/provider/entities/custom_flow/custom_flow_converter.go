package custom_flow

import (
	sdkCustomFlow "github.com/control-monkey/controlmonkey-sdk-go/services/custom_flow"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*sdkCustomFlow.CustomFlow, bool) {
	var retVal *sdkCustomFlow.CustomFlow

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(sdkCustomFlow.CustomFlow)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if plan.Name != state.Name {
		retVal.SetName(plan.Name.ValueStringPointer())
		hasChanges = true
	}
	if plan.FlowYaml != state.FlowYaml {
		retVal.SetFlowYaml(plan.FlowYaml.ValueStringPointer())
		hasChanges = true
	}
	if plan.IsEnabled != state.IsEnabled {
		retVal.SetIsEnabled(plan.IsEnabled.ValueBoolPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}
