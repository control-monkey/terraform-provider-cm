package customAbacConfiguration

import (
	apiCustomAbacConfiguration "github.com/control-monkey/controlmonkey-sdk-go/services/custom_abac_configuration"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res *apiCustomAbacConfiguration.CustomAbacConfiguration, state *ResourceModel) {
	state.Name = helpers.StringValueOrNull(res.Name)
	state.CustomAbacId = helpers.StringValueOrNull(res.CustomAbacId)

	if res.Roles != nil {
		ec := updateStateAfterReadRoles(res.Roles)
		state.Roles = ec
	} else {
		state.Roles = nil
	}
}

func updateStateAfterReadRoles(elements []*apiCustomAbacConfiguration.Role) []*RoleModel {
	var retVal []*RoleModel

	if elements != nil {
		retVal = make([]*RoleModel, 0)

		for _, rule := range elements {
			cp := updateStateAfterReadRole(rule)
			retVal = append(retVal, &cp)
		}
	}

	return retVal
}

func updateStateAfterReadRole(element *apiCustomAbacConfiguration.Role) RoleModel {
	var retVal RoleModel

	retVal.OrgId = helpers.StringValueOrNull(element.OrgId)
	retVal.OrgRole = helpers.StringValueOrNull(element.OrgRole)
	retVal.TeamIds = helpers.StringPointerSliceToTfList(element.TeamIds)

	return retVal
}
