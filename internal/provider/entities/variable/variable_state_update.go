package variable

import (
	"fmt"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	sdkVariable "github.com/control-monkey/controlmonkey-sdk-go/services/variable"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(res *sdkVariable.ReadVariableOutput, state *ResourceModel) {
	variable := res.Variable

	state.Scope = helpers.StringValueOrNull(variable.Scope)
	state.ScopeId = helpers.StringValueOrNull(variable.ScopeId)
	state.Key = helpers.StringValueOrNull(variable.Key)
	state.Type = helpers.StringValueOrNull(variable.Type)

	// if it's sensitive, we take the value from the state file because the api does not respond secret values.
	// if it's not sensitive, we take the value from the response.
	if state.IsSensitive.ValueBool() == false {
		state.Value = helpers.StringValueOrNull(variable.Value)
	}

	state.DisplayName = helpers.StringValueOrNull(variable.DisplayName)
	state.IsSensitive = helpers.BoolValueOrNull(variable.IsSensitive)
	state.IsOverridable = helpers.BoolValueOrNull(variable.IsOverridable)
	state.IsRequired = helpers.BoolValueOrNull(variable.IsRequired)
	state.Description = helpers.StringValueIfNotEqual(variable.Description, "")

	if variable.ValueConditions != nil {
		vc := updateStateAfterReadValueConditions(variable.ValueConditions)
		state.ValueConditions = vc
	} else {
		state.ValueConditions = nil
	}
}

func updateStateAfterReadValueConditions(valueConditions []*sdkVariable.Condition) []*ConditionModel {
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

func updateStateAfterReadCondition(condition *sdkVariable.Condition) ConditionModel {
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
