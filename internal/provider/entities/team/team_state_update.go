package team

import (
	sdkTeam "github.com/control-monkey/controlmonkey-sdk-go/services/team"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res *sdkTeam.Team, state *ResourceModel) {
	team := res

	state.Name = helpers.StringValueOrNull(team.Name)
	state.CustomIdpId = helpers.StringValueOrNull(team.CustomIdpId)
}
