package controlPolicyGroupMappings

import (
	sdkControlPolicyGroup "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy_group"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(apiEntities []*sdkControlPolicyGroup.ControlPolicyGroupMapping, state *ResourceModel) {
	state.ControlPolicyGroupId = state.ID

	if apiEntities != nil {
		targets := updateStateAfterReadTargets(apiEntities)
		state.Targets = targets
	} else {
		state.Targets = nil
	}
}

func updateStateAfterReadTargets(apiEntities []*sdkControlPolicyGroup.ControlPolicyGroupMapping) []*TargetModel {
	var retVal []*TargetModel

	if len(apiEntities) > 0 {
		retVal = make([]*TargetModel, len(apiEntities))

		for i, target := range apiEntities {
			tfEntity := updateStateAfterReadTarget(target)
			retVal[i] = &tfEntity
		}
	}

	return retVal
}

func updateStateAfterReadTarget(apiEntity *sdkControlPolicyGroup.ControlPolicyGroupMapping) TargetModel {
	var retVal TargetModel

	retVal.TargetId = helpers.StringValueOrNull(apiEntity.TargetId)
	retVal.TargetType = helpers.StringValueOrNull(apiEntity.TargetType)
	retVal.EnforcementLevel = helpers.StringValueOrNull(apiEntity.EnforcementLevel)

	if apiEntity.OverrideEnforcements != nil {
		tfOverrideEnforcement := updateStateAfterReadOverrideEnforcements(apiEntity.OverrideEnforcements)
		retVal.OverrideEnforcements = tfOverrideEnforcement
	} else {
		retVal.OverrideEnforcements = nil
	}

	return retVal
}

func updateStateAfterReadOverrideEnforcements(apiEntities []*sdkControlPolicyGroup.OverrideEnforcement) []*OverrideEnforcementModel {
	var retVal []*OverrideEnforcementModel

	if len(apiEntities) > 0 {
		retVal = make([]*OverrideEnforcementModel, len(apiEntities))

		for i, apiEntity := range apiEntities {
			tfEntity := updateStateAfterReadOverrideEnforcement(apiEntity)
			retVal[i] = &tfEntity
		}
	}

	return retVal
}

func updateStateAfterReadOverrideEnforcement(apiEntity *sdkControlPolicyGroup.OverrideEnforcement) OverrideEnforcementModel {
	var retVal OverrideEnforcementModel

	retVal.ControlPolicyId = helpers.StringValueOrNull(apiEntity.ControlPolicyId)
	retVal.EnforcementLevel = helpers.StringValueOrNull(apiEntity.EnforcementLevel)
	retVal.StackIds = helpers.StringPointerSliceToTfList(apiEntity.StackIds)

	return retVal
}
