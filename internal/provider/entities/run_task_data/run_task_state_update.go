package run_task_data

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkRunTask "github.com/control-monkey/controlmonkey-sdk-go/services/run_task"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(apiEntity *sdkRunTask.RunTask, state *ResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(controlmonkey.StringValue(apiEntity.ID))
	state.Name = types.StringValue(controlmonkey.StringValue(apiEntity.Name))
}
