package run_task

import (
	sdkRunTask "github.com/control-monkey/controlmonkey-sdk-go/services/run_task"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res *sdkRunTask.RunTask, state *ResourceModel) {
	state.Name = helpers.StringValueOrNull(res.Name)
	state.Url = helpers.StringValueOrNull(res.Url)
	state.IsEnabled = helpers.BoolValueOrNull(res.IsEnabled)
	state.IsHmacKeyConfigured = helpers.BoolValueOrNull(res.IsHmacKeyConfigured)
}
