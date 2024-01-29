package variable

import (
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/controlmonkey-sdk-go/services/variable"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"reflect"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*variable.Variable, bool) {
	var retVal *variable.Variable

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(variable.Variable)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if plan.Scope != state.Scope {
		retVal.SetScope(plan.Scope.ValueStringPointer())
		hasChanges = true
	}
	if plan.ScopeId != state.ScopeId {
		retVal.SetScopeId(plan.ScopeId.ValueStringPointer())
		hasChanges = true
	}
	if plan.Key != state.Key {
		retVal.SetKey(plan.Key.ValueStringPointer())
		hasChanges = true
	}
	if plan.Type != state.Type {
		retVal.SetType(plan.Type.ValueStringPointer())
		hasChanges = true
	}
	if plan.Value != state.Value {
		retVal.SetValue(plan.Value.ValueStringPointer())
		hasChanges = true
	}
	if plan.IsSensitive != state.IsSensitive {
		retVal.SetIsSensitive(plan.IsSensitive.ValueBoolPointer())
		hasChanges = true
	}
	if plan.IsOverridable != state.IsOverridable {
		retVal.SetIsOverridable(plan.IsOverridable.ValueBoolPointer())
		hasChanges = true
	}
	if plan.IsRequired != state.IsRequired {
		retVal.SetIsRequired(plan.IsRequired.ValueBoolPointer())
		hasChanges = true
	}
	if plan.Description != state.Description {
		retVal.SetDescription(plan.Description.ValueStringPointer())
		hasChanges = true
	}

	if vc, hasChanged := valueConditionsConverter(plan.ValueConditions, state.ValueConditions, converterType); hasChanged {
		retVal.SetValueConditions(vc)
		hasChanges = true
	}

	return retVal, hasChanges
}

func valueConditionsConverter(plan []*ConditionModel, state []*ConditionModel, converterType commons.ConverterType) ([]*variable.Condition, bool) {
	var retVal []*variable.Condition
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true
		retVal = make([]*variable.Condition, 0)

		for _, r := range plan {
			rule := conditionConverter(r)
			retVal = append(retVal, rule)
		}
	}

	return retVal, hasChanged
}

func conditionConverter(plan *ConditionModel) *variable.Condition {
	retVal := new(variable.Condition)

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
