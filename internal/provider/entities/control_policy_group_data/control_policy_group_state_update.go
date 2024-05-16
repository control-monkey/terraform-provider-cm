package control_policy_group_data

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkControlPolicyGroup "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy_group"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(apiEntity *sdkControlPolicyGroup.ControlPolicyGroup, state *ResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(controlmonkey.StringValue(apiEntity.ID))
	state.Name = types.StringValue(controlmonkey.StringValue(apiEntity.Name))
}
