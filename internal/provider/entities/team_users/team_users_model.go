package team_users

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/team"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID     types.String `tfsdk:"id"`
	TeamId types.String `tfsdk:"team_id"`
	Users  []*UserModel `tfsdk:"users"`
}

type UserModel struct { //When new field is added consider Hash() function
	Email types.String `tfsdk:"email"`
}

type MergedEntities struct {
	EntitiesToCreate []*team.TeamUser
	EntitiesToDelete []*team.TeamUser
}

func (e *UserModel) Hash() string {
	return e.Email.ValueString()
}

func (e *UserModel) GetBlockIdentifier() string {
	retVal := ""

	if helpers.IsKnown(e.Email) {
		retVal += e.Hash() // do not use e.Hash if another property is added to Model
	}

	return retVal
}
