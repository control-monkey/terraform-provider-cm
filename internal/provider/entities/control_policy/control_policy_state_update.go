package controlPolicy

import (
	"encoding/json"
	apiControlPolicy "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
)

func UpdateStateAfterRead(res *apiControlPolicy.ControlPolicy, state *ResourceModel) {
	state.Name = helpers.StringValueOrNull(res.Name)
	state.Description = helpers.StringValueIfNotEqual(res.Description, "")
	state.Type = helpers.StringValueOrNull(res.Type)

	jsonSettingsString, err := json.Marshal(res.Parameters)
	if err != nil {
		state.Parameters = jsontypes.NewNormalizedNull()
	}

	state.Parameters = jsontypes.NewNormalizedValue(string(jsonSettingsString))
}
