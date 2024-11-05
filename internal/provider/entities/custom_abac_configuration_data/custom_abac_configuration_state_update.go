package custom_abac_configuration_data

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkCustomAbacConfiguration "github.com/control-monkey/controlmonkey-sdk-go/services/custom_abac_configuration"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(apiEntity *sdkCustomAbacConfiguration.CustomAbacConfiguration, state *ResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(controlmonkey.StringValue(apiEntity.ID))
	state.Name = types.StringValue(controlmonkey.StringValue(apiEntity.Name))
}
