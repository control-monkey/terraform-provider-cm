package customFlowMappings

import (
	sdkCustomFlow "github.com/control-monkey/controlmonkey-sdk-go/services/custom_flow"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(apiEntities []*sdkCustomFlow.CustomFlowMapping, state *ResourceModel) {
	state.CustomFlowId = state.ID

	// The API returns every mapping in the organization, so keep only the ones this resource owns.
	owned := filterByCustomFlowId(apiEntities, state.ID.ValueString())

	if owned != nil {
		targets := updateStateAfterReadTargets(owned)
		state.Targets = targets
	} else {
		state.Targets = nil
	}
}

func filterByCustomFlowId(apiEntities []*sdkCustomFlow.CustomFlowMapping, customFlowId string) []*sdkCustomFlow.CustomFlowMapping {
	if apiEntities == nil {
		return nil
	}

	retVal := make([]*sdkCustomFlow.CustomFlowMapping, 0, len(apiEntities))

	for _, e := range apiEntities {
		if e != nil && e.CustomFlowId != nil && *e.CustomFlowId == customFlowId {
			retVal = append(retVal, e)
		}
	}

	return retVal
}

func updateStateAfterReadTargets(targets []*sdkCustomFlow.CustomFlowMapping) []*TargetModel {
	var retVal []*TargetModel

	if len(targets) > 0 {
		retVal = make([]*TargetModel, len(targets))

		for i, target := range targets {
			u := updateStateAfterReadTarget(target)
			retVal[i] = &u
		}
	}

	return retVal
}

func updateStateAfterReadTarget(target *sdkCustomFlow.CustomFlowMapping) TargetModel {
	var retVal TargetModel

	retVal.TargetId = helpers.StringValueOrNull(target.TargetId)
	retVal.TargetType = helpers.StringValueOrNull(target.TargetType)

	return retVal
}
