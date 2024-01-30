package controlPolicyMapping

import (
	sdkControlPolicy "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res *sdkControlPolicy.ControlPolicyMapping, state *ResourceModel) {
	state.ControlPolicyId = helpers.StringValueOrNull(res.ControlPolicyId)
	state.TargetId = helpers.StringValueOrNull(res.TargetId)
	state.TargetType = helpers.StringValueOrNull(res.TargetType)
	state.EnforcementLevel = helpers.StringValueOrNull(res.EnforcementLevel)
	state.ID = ComputeId(res)
}
