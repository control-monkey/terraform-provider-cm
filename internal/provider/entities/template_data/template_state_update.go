package template_data

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkTemplate "github.com/control-monkey/controlmonkey-sdk-go/services/template"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(apiEntity *sdkTemplate.Template, state *ResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(controlmonkey.StringValue(apiEntity.ID))
	state.Name = types.StringValue(controlmonkey.StringValue(apiEntity.Name))
}
