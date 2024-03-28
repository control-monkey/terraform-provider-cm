package controlPolicy

import (
	apiControlPolicy "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*apiControlPolicy.ControlPolicy, bool) {
	var retVal *apiControlPolicy.ControlPolicy

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(apiControlPolicy.ControlPolicy)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if plan.Name != state.Name {
		retVal.SetName(plan.Name.ValueStringPointer())
		hasChanges = true
	}
	if plan.Description != state.Description {
		retVal.SetDescription(plan.Description.ValueStringPointer())
		hasChanges = true
	}
	if plan.Type != state.Type {
		retVal.SetType(plan.Type.ValueStringPointer())
		hasChanges = true
	}
	if plan.Parameters != state.Parameters {
		a := new(map[string]any)
		plan.Parameters.Unmarshal(a)
		retVal.SetParameters(a)

		if retVal.Type == nil { //if parameters was changed, type must be sent
			retVal.SetType(plan.Type.ValueStringPointer())
		}

		hasChanges = true
	}

	return retVal, hasChanges
}
