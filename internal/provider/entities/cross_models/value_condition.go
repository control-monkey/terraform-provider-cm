package cross_models

import (
	"fmt"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	apiCrossModels "github.com/control-monkey/controlmonkey-sdk-go/services/cross_models"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
)

//region Model

type ConditionModel struct {
	Operator types.String `tfsdk:"operator"`
	Value    types.String `tfsdk:"value"`
	Values   types.List   `tfsdk:"values"`
}

//endregion

//region Create/Update Converter

func ValueConditionsConverter(plan []*ConditionModel, state []*ConditionModel, converterType commons.ConverterType) ([]*apiCrossModels.Condition, bool) {
	var retVal []*apiCrossModels.Condition
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true

		if plan != nil {
			retVal = make([]*apiCrossModels.Condition, 0)

			for _, r := range plan {
				rule := conditionConverter(r)
				retVal = append(retVal, rule)
			}
		}
	}

	return retVal, hasChanged
}

func conditionConverter(plan *ConditionModel) *apiCrossModels.Condition {
	retVal := new(apiCrossModels.Condition)

	operator := plan.Operator
	retVal.SetOperator(operator.ValueStringPointer())

	// We rely on the assumption that ValueString() is used without checking the pointer only when it must appear.
	planValue := plan.Value

	switch op := operator.ValueString(); op {
	case cmTypes.Ne:
		var strVal any = planValue.ValueString()
		retVal.SetValue(&strVal)
	case cmTypes.Gt, cmTypes.Gte, cmTypes.Lt, cmTypes.Lte:
		var intVal any
		_, num := helpers.CheckAndGetIfNumericString(planValue.ValueString()) // was already checked that it is numeric
		intVal = num
		retVal.SetValue(&intVal)
	case cmTypes.In:
		var sliceVal any = helpers.TfListToStringPointerSlice(plan.Values)
		retVal.SetValue(&sliceVal)
	case cmTypes.StartsWith, cmTypes.Contains:
		var strVal any = planValue.ValueString()
		retVal.SetValue(&strVal)
	}

	return retVal
}

//endregion

//region Update State After Read

func UpdateStateAfterReadValueConditions(valueConditions []*apiCrossModels.Condition) []*ConditionModel {
	var retVal []*ConditionModel

	if valueConditions != nil {
		retVal = make([]*ConditionModel, 0)

		for _, condition := range valueConditions {
			c := updateStateAfterReadCondition(condition)
			retVal = append(retVal, &c)
		}
	}

	return retVal
}

func updateStateAfterReadCondition(condition *apiCrossModels.Condition) ConditionModel {
	var retVal ConditionModel

	operator := condition.Operator
	retVal.Operator = helpers.StringValueOrNull(operator)
	retVal.Values = types.ListNull(types.StringType) // set default null of type string. otherwise missing type error occurs.

	switch op := *operator; op {
	case cmTypes.Ne:
		var strValue = (*condition.Value).(string)
		retVal.Value = helpers.StringValueOrNull(&strValue)
	case cmTypes.Gt, cmTypes.Gte, cmTypes.Lt, cmTypes.Lte:
		var floatVal = (*condition.Value).(float64)
		strVal := fmt.Sprint(floatVal)
		retVal.Value = helpers.StringValueOrNull(&strVal)
	case cmTypes.In:
		retVal.Values = helpers.StringPointerSliceToTfList(condition.Values) //Note Values, not Value
	case cmTypes.StartsWith, cmTypes.Contains:
		var strValue = (*condition.Value).(string)
		retVal.Value = helpers.StringValueOrNull(&strValue)
	}

	return retVal
}

//endregion
