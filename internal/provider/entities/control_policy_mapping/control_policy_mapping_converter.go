package controlPolicyMapping

import (
	controlPolicy "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*controlPolicy.ControlPolicyMapping, bool) {
	var retVal *controlPolicy.ControlPolicyMapping

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(controlPolicy.ControlPolicyMapping)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	//The following fields together are the identifier of the resource
	retVal.SetControlPolicyId(plan.ControlPolicyId.ValueStringPointer())
	retVal.SetTargetId(plan.TargetId.ValueStringPointer())
	retVal.SetTargetType(plan.TargetType.ValueStringPointer())

	if plan.EnforcementLevel != state.EnforcementLevel {
		retVal.SetEnforcementLevel(plan.EnforcementLevel.ValueStringPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}
