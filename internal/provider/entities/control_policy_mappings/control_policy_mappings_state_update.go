package controlPolicyMappings

import (
	sdkControlPolicy "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(apiEntities []*sdkControlPolicy.ControlPolicyMapping, state *ResourceModel) {
	state.ControlPolicyId = state.ID

	if apiEntities != nil {
		targets := updateStateAfterReadTargets(apiEntities)
		state.Targets = targets
	} else {
		state.Targets = nil
	}
}

func updateStateAfterReadTargets(targets []*sdkControlPolicy.ControlPolicyMapping) []*TargetModel {
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

func updateStateAfterReadTarget(target *sdkControlPolicy.ControlPolicyMapping) TargetModel {
	var retVal TargetModel

	retVal.TargetId = helpers.StringValueOrNull(target.TargetId)
	retVal.TargetType = helpers.StringValueOrNull(target.TargetType)
	retVal.EnforcementLevel = helpers.StringValueOrNull(target.EnforcementLevel)

	return retVal
}
