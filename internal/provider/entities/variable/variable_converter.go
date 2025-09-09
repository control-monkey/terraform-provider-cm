package variable

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/variable"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
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
	if plan.DisplayName != state.DisplayName {
		retVal.SetDisplayName(plan.DisplayName.ValueStringPointer())
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

	if vc, hasChanged := cross_models.ValueConditionsConverter(plan.ValueConditions, state.ValueConditions, converterType); hasChanged {
		retVal.SetValueConditions(vc)
		hasChanges = true
	}

	if plan.BlueprintVariableManagedBy != state.BlueprintVariableManagedBy {
		retVal.SetBlueprintManagedBy(plan.BlueprintVariableManagedBy.ValueStringPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}
