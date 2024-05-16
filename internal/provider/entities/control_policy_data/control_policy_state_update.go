package control_policy_data

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkControlPolicy "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(apiEntity *sdkControlPolicy.ControlPolicy, state *ResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(controlmonkey.StringValue(apiEntity.ID))
	state.Name = types.StringValue(controlmonkey.StringValue(apiEntity.Name))
}
