package custom_role_data

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkCustomRole "github.com/control-monkey/controlmonkey-sdk-go/services/custom_role"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(apiEntity *sdkCustomRole.CustomRole, state *ResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(controlmonkey.StringValue(apiEntity.ID))
	state.Name = types.StringValue(controlmonkey.StringValue(apiEntity.Name))
}
