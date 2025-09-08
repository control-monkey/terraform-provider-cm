package cross_models

import (
	sdkCrossModels "github.com/control-monkey/controlmonkey-sdk-go/services/cross_models"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//region Model

type RunnerConfigModel struct {
	Mode   types.String `tfsdk:"mode"`
	Groups types.List   `tfsdk:"groups"`
}

//endregion

//region Create/Update Converter

func RunnerConfigConverter(plan *RunnerConfigModel, state *RunnerConfigModel, converterType commons.ConverterType) (*sdkCrossModels.RunnerConfig, bool) {
	var retVal *sdkCrossModels.RunnerConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(sdkCrossModels.RunnerConfig)
	hasChanges := false

	if state == nil {
		state = new(RunnerConfigModel) // dummy initialization
		hasChanges = true              // must have changes because before is null and after is not
	}

	if plan.Mode != state.Mode {
		retVal.SetMode(plan.Mode.ValueStringPointer())
		hasChanges = true
	}

	if innerProperty, hasInnerChanges := helpers.TfListStringConverter(plan.Groups, state.Groups); hasInnerChanges {
		retVal.SetGroups(innerProperty)
		hasChanges = true
	}

	return retVal, hasChanges
}

//endregion

//region Update State After Read

func UpdateStateAfterReadRunnerConfig(rc *sdkCrossModels.RunnerConfig) RunnerConfigModel {
	var retVal RunnerConfigModel

	if rc != nil {
		retVal.Mode = helpers.StringValueOrNull(rc.Mode)
		retVal.Groups = helpers.StringPointerSliceToTfList(rc.Groups)
	}

	return retVal
}

//endregion
