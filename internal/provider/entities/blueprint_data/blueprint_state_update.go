package blueprint_data

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkBlueprint "github.com/control-monkey/controlmonkey-sdk-go/services/blueprint"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(apiEntity *sdkBlueprint.Blueprint, state *ResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(controlmonkey.StringValue(apiEntity.ID))
	state.Name = types.StringValue(controlmonkey.StringValue(apiEntity.Name))
}
