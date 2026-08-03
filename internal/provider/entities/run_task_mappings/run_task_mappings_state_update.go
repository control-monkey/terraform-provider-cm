package runTaskMappings

import (
	sdkRunTask "github.com/control-monkey/controlmonkey-sdk-go/services/run_task"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(apiEntities []*sdkRunTask.RunTaskMapping, state *ResourceModel) {
	state.RunTaskId = state.ID

	if apiEntities != nil {
		targets := updateStateAfterReadTargets(apiEntities)
		state.Targets = targets
	} else {
		state.Targets = nil
	}
}

func updateStateAfterReadTargets(targets []*sdkRunTask.RunTaskMapping) []*TargetModel {
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

func updateStateAfterReadTarget(target *sdkRunTask.RunTaskMapping) TargetModel {
	var retVal TargetModel

	retVal.TargetId = helpers.StringValueOrNull(target.TargetId)
	retVal.TargetType = helpers.StringValueOrNull(target.TargetType)
	retVal.EnforcementLevel = helpers.StringValueOrNull(target.EnforcementLevel)
	retVal.Stage = helpers.StringValueOrNull(target.Stage)

	return retVal
}
