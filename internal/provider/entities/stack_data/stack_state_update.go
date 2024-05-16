package stack_data

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkStack "github.com/control-monkey/controlmonkey-sdk-go/services/stack"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(apiEntity *sdkStack.Stack, state *ResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(controlmonkey.StringValue(apiEntity.ID))
	state.Name = types.StringValue(controlmonkey.StringValue(apiEntity.Name))
	state.NamespaceId = types.StringValue(controlmonkey.StringValue(apiEntity.NamespaceId))
}
