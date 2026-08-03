package custom_flow

import (
	sdkCustomFlow "github.com/control-monkey/controlmonkey-sdk-go/services/custom_flow"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res *sdkCustomFlow.CustomFlow, state *ResourceModel) {
	state.Name = helpers.StringValueOrNull(res.Name)
	state.FlowYaml = helpers.StringValueOrNull(res.FlowYaml)
	state.IsEnabled = helpers.BoolValueOrNull(res.IsEnabled)
}
