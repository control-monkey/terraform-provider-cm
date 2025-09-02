package cross_models

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/cross_models"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//region Model

type AutoSyncModel struct {
	DeployWhenDriftDetected types.Bool `tfsdk:"deploy_when_drift_detected"`
}

//endregion

//region Create/Update Converter

func AutoSyncConverter(plan *AutoSyncModel, state *AutoSyncModel, converterType commons.ConverterType) (*cross_models.AutoSync, bool) {
	var retVal *cross_models.AutoSync

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(cross_models.AutoSync)
	hasChanges := false

	if state == nil {
		state = new(AutoSyncModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if plan.DeployWhenDriftDetected != state.DeployWhenDriftDetected {
		retVal.SetDeployWhenDriftDetected(plan.DeployWhenDriftDetected.ValueBoolPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}

//endregion

//region Update State After Read

func UpdateStateAfterReadAutoSync(as *cross_models.AutoSync) AutoSyncModel {
	var retVal AutoSyncModel

	retVal.DeployWhenDriftDetected = helpers.BoolValueOrNull(as.DeployWhenDriftDetected)

	return retVal
}

//endregion
