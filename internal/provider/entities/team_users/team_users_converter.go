package team_users

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/team"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/hashicorp/go-set/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Merge(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) *MergedEntities {
	retVal := new(MergedEntities)

	if plan == nil {
		plan = new(ResourceModel) // delete merger
	}

	if state == nil {
		state = new(ResourceModel) // create merger
	}

	var teamId types.String
	if plan.TeamId.IsNull() == false {
		teamId = plan.TeamId
	} else {
		teamId = state.TeamId
	}

	planUsers := set.HashSetFrom[*UserModel, string](plan.Users)
	stateUsers := set.HashSetFrom[*UserModel, string](state.Users)

	usersToAdd := planUsers.Difference(stateUsers)
	usersToDelete := stateUsers.Difference(planUsers)

	retVal.EntitiesToCreate = buildEntities(usersToAdd, teamId)
	retVal.EntitiesToDelete = buildEntities(usersToDelete, teamId)

	return retVal
}

func buildEntities(entities set.Collection[*UserModel], teamId types.String) []*team.TeamUser {
	retVal := make([]*team.TeamUser, entities.Size())

	for i, u := range entities.Slice() {
		tu := new(team.TeamUser)
		tu.SetTeamId(teamId.ValueStringPointer())
		tu.SetUserEmail(u.Email.ValueStringPointer())

		retVal[i] = tu
	}

	return retVal
}
