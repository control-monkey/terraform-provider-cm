package cross_models

import (
	sdkCrossModels "github.com/control-monkey/controlmonkey-sdk-go/services/cross_models"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//region Model

type DeploymentBehaviorModel struct {
	DeployOnPush    types.Bool `tfsdk:"deploy_on_push"`
	WaitForApproval types.Bool `tfsdk:"wait_for_approval"`
}

//endregion

//region Create/Update Converter

func DeploymentBehaviorConverter(plan *DeploymentBehaviorModel, state *DeploymentBehaviorModel, converterType commons.ConverterType) (*sdkCrossModels.DeploymentBehavior, bool) {
	var retVal *sdkCrossModels.DeploymentBehavior

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(sdkCrossModels.DeploymentBehavior)
	hasChanges := false

	if state == nil {
		state = new(DeploymentBehaviorModel) // dummy initialization
		hasChanges = true                    // must have changes because before is null and after is not
	}

	if plan.DeployOnPush != state.DeployOnPush {
		retVal.SetDeployOnPush(plan.DeployOnPush.ValueBoolPointer())
		hasChanges = true
	}
	if plan.WaitForApproval != state.WaitForApproval {
		retVal.SetWaitForApproval(plan.WaitForApproval.ValueBoolPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}

//endregion

//region Update State After Read

func UpdateStateAfterReadDeploymentBehavior(deploymentBehavior *sdkCrossModels.DeploymentBehavior) DeploymentBehaviorModel {
	var retVal DeploymentBehaviorModel

	retVal.DeployOnPush = helpers.BoolValueOrNull(deploymentBehavior.DeployOnPush)
	retVal.WaitForApproval = helpers.BoolValueOrNull(deploymentBehavior.WaitForApproval)

	return retVal
}

//endregion
