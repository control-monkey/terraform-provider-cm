package namespace_data

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkNamespace "github.com/control-monkey/controlmonkey-sdk-go/services/namespace"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(apiEntity *sdkNamespace.Namespace, state *ResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(controlmonkey.StringValue(apiEntity.ID))
	state.Name = types.StringValue(controlmonkey.StringValue(apiEntity.Name))
}
