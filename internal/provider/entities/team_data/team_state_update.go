package team_data

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkTeam "github.com/control-monkey/controlmonkey-sdk-go/services/team"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(apiEntity *sdkTeam.Team, state *ResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(controlmonkey.StringValue(apiEntity.ID))
	state.Name = types.StringValue(controlmonkey.StringValue(apiEntity.Name))
}
