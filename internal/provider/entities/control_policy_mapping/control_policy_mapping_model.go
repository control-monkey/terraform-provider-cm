package controlPolicyMapping

import (
	"fmt"
	sdkControlPolicy "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID               types.String `tfsdk:"id"`
	ControlPolicyId  types.String `tfsdk:"control_policy_id"`
	TargetId         types.String `tfsdk:"target_id"`
	TargetType       types.String `tfsdk:"target_type"`
	EnforcementLevel types.String `tfsdk:"enforcement_level"`
}

func ComputeId(mapping *sdkControlPolicy.ControlPolicyMapping) types.String {
	var retVal types.String

	id := fmt.Sprintf("%s/%s/%s", *mapping.ControlPolicyId, *mapping.TargetId, *mapping.TargetType)
	retVal = helpers.StringValueOrNull(&id)

	return retVal
}
