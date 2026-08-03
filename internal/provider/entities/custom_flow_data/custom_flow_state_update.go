package custom_flow_data

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkCustomFlow "github.com/control-monkey/controlmonkey-sdk-go/services/custom_flow"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(apiEntity *sdkCustomFlow.CustomFlow, state *ResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(controlmonkey.StringValue(apiEntity.ID))
	state.Name = types.StringValue(controlmonkey.StringValue(apiEntity.Name))
}
