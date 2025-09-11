package cross_models

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/cross_models"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//region Model

type RunTriggerModel struct {
	Patterns        types.List `tfsdk:"patterns"`
	ExcludePatterns types.List `tfsdk:"exclude_patterns"`
}

//endregion

//region Create/Update Converter

func RunTriggerConverter(plan *RunTriggerModel, state *RunTriggerModel, converterType commons.ConverterType) (*cross_models.RunTrigger, bool) {
	var retVal *cross_models.RunTrigger

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(cross_models.RunTrigger)
	hasChanges := false

	if state == nil {
		state = new(RunTriggerModel) // dummy initialization
		hasChanges = true            // must have changes because before is null and after is not
	}

	if innerProperty, hasInnerChanges := helpers.TfListStringConverter(plan.Patterns, state.Patterns); hasInnerChanges {
		retVal.SetPatterns(innerProperty)
		hasChanges = true
	}

	if innerProperty, hasInnerChanges := helpers.TfListStringConverter(plan.ExcludePatterns, state.ExcludePatterns); hasInnerChanges {
		retVal.SetExcludePatterns(innerProperty)
		hasChanges = true
	}

	return retVal, hasChanges
}

//endregion

//region Update State After Read

func UpdateStateAfterReadRunTrigger(runTrigger *cross_models.RunTrigger) RunTriggerModel {
	var retVal RunTriggerModel

	retVal.Patterns = helpers.StringPointerSliceToTfList(runTrigger.Patterns)
	retVal.ExcludePatterns = helpers.StringPointerSliceToTfList(runTrigger.ExcludePatterns)

	return retVal
}

//endregion
