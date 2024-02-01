package team_users

import (
	sdkTeam "github.com/control-monkey/controlmonkey-sdk-go/services/team"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res []*sdkTeam.TeamUser, state *ResourceModel) {
	teamUsers := res

	state.TeamId = state.ID

	if teamUsers != nil {
		ec := updateStateAfterReadTeamUsers(teamUsers)
		state.Users = ec
	} else {
		state.Users = nil
	}
}

func updateStateAfterReadTeamUsers(teamUsers []*sdkTeam.TeamUser) []*UserModel {
	var retVal []*UserModel

	if len(teamUsers) > 0 {
		retVal = make([]*UserModel, len(teamUsers))

		for i, user := range teamUsers {
			u := updateStateAfterReadUser(user)
			retVal[i] = &u
		}
	}

	return retVal
}

func updateStateAfterReadUser(user *sdkTeam.TeamUser) UserModel {
	var retVal UserModel

	retVal.Email = helpers.StringValueOrNull(user.UserEmail)

	return retVal
}
