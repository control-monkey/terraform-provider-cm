package controlPolicyGroup

import (
	apiControlPolicyGroup "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy_group"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res *apiControlPolicyGroup.ControlPolicyGroup, state *ResourceModel) {
	state.Name = helpers.StringValueOrNull(res.Name)
	state.Description = helpers.StringValueIfNotEqual(res.Description, "")

	if res.ControlPolicies != nil {
		ec := updateStateAfterReadControlPolicies(res.ControlPolicies)
		state.ControlPolicies = ec
	} else {
		state.ControlPolicies = nil
	}
}
func updateStateAfterReadControlPolicies(controlPolicies []*apiControlPolicyGroup.ControlPolicy) []*ControlPolicyModel {
	var retVal []*ControlPolicyModel

	if controlPolicies != nil {
		retVal = make([]*ControlPolicyModel, 0)

		for _, rule := range controlPolicies {
			cp := updateStateAfterReadControlPolicy(rule)
			retVal = append(retVal, &cp)
		}
	}

	return retVal
}

func updateStateAfterReadControlPolicy(credentials *apiControlPolicyGroup.ControlPolicy) ControlPolicyModel {
	var retVal ControlPolicyModel

	retVal.ControlPolicyId = helpers.StringValueOrNull(credentials.ControlPolicyId)
	retVal.Severity = helpers.StringValueOrNull(credentials.Severity)

	return retVal
}
