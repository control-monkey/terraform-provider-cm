package variable

import (
	sdkVariable "github.com/control-monkey/controlmonkey-sdk-go/services/variable"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
)

func UpdateStateAfterRead(res *sdkVariable.ReadVariableOutput, state *ResourceModel) {
	variable := res.Variable

	state.Scope = helpers.StringValueOrNull(variable.Scope)
	state.ScopeId = helpers.StringValueOrNull(variable.ScopeId)
	state.Key = helpers.StringValueOrNull(variable.Key)
	state.Type = helpers.StringValueOrNull(variable.Type)

	// if it's sensitive, we take the value from the state file because the api does not respond secret values.
	// if it's not sensitive, we take the value from the response.
	if state.IsSensitive.ValueBool() == false {
		state.Value = helpers.StringValueOrNull(variable.Value)
	}

	state.DisplayName = helpers.StringValueOrNull(variable.DisplayName)
	state.IsSensitive = helpers.BoolValueOrNull(variable.IsSensitive)
	state.IsOverridable = helpers.BoolValueOrNull(variable.IsOverridable)
	state.IsRequired = helpers.BoolValueOrNull(variable.IsRequired)
	state.Description = helpers.StringValueIfNotEqual(variable.Description, "")

	if variable.ValueConditions != nil {
		vc := cross_models.UpdateStateAfterReadValueConditions(variable.ValueConditions)
		state.ValueConditions = vc
	} else {
		state.ValueConditions = nil
	}
}
