package controlPolicyGroup

import (
	apiControlPolicyGroup "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy_group"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"reflect"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*apiControlPolicyGroup.ControlPolicyGroup, bool) {
	var retVal *apiControlPolicyGroup.ControlPolicyGroup

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(apiControlPolicyGroup.ControlPolicyGroup)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if plan.Name != state.Name {
		retVal.SetName(plan.Name.ValueStringPointer())
		hasChanges = true
	}
	if plan.Description != state.Description {
		retVal.SetDescription(plan.Description.ValueStringPointer())
		hasChanges = true
	}
	if cp, hasChanged := controlPoliciesConverter(plan.ControlPolicies, state.ControlPolicies, converterType); hasChanged {
		retVal.SetControlPolicies(cp)
		hasChanges = true
	}
	return retVal, hasChanges
}

func controlPoliciesConverter(plan []*ControlPolicyModel, state []*ControlPolicyModel, converterType commons.ConverterType) ([]*apiControlPolicyGroup.ControlPolicy, bool) {
	var retVal []*apiControlPolicyGroup.ControlPolicy
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true

		if plan != nil {
			retVal = make([]*apiControlPolicyGroup.ControlPolicy, 0)

			for _, r := range plan {
				policy := controlPolicyConverter(r)
				retVal = append(retVal, policy)
			}
		}
	}

	return retVal, hasChanged
}

func controlPolicyConverter(plan *ControlPolicyModel) *apiControlPolicyGroup.ControlPolicy {
	retVal := new(apiControlPolicyGroup.ControlPolicy)

	retVal.SetControlPolicyId(plan.ControlPolicyId.ValueStringPointer())
	retVal.SetSeverity(plan.Severity.ValueStringPointer())

	return retVal
}
