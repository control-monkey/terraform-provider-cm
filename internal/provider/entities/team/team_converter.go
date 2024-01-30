package team

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/team"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*team.Team, bool) {
	var retVal *team.Team

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(team.Team)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if plan.Name != state.Name {
		retVal.SetName(plan.Name.ValueStringPointer())
		hasChanges = true
	}
	if plan.CustomIdpId != state.CustomIdpId {
		retVal.SetCustomIdpId(plan.CustomIdpId.ValueStringPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}
