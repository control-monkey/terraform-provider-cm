package customRole

import (
	apiCustomRole "github.com/control-monkey/controlmonkey-sdk-go/services/custom_role"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res *apiCustomRole.CustomRole, state *ResourceModel) {
	state.Name = helpers.StringValueOrNull(res.Name)
	state.Description = helpers.StringValueIfNotEqual(res.Description, "")

	if res.Permissions != nil {
		ec := updateStateAfterReadPermissions(res.Permissions)
		state.Permissions = ec
	} else {
		state.Permissions = nil
	}

	state.StackRestriction = helpers.StringValueOrNull(res.StackRestriction)
}
func updateStateAfterReadPermissions(permissions []*apiCustomRole.Permission) []*PermissionModel {
	var retVal []*PermissionModel

	if permissions != nil {
		retVal = make([]*PermissionModel, 0)

		for _, rule := range permissions {
			cp := updateStateAfterReadPermission(rule)
			retVal = append(retVal, &cp)
		}
	}

	return retVal
}

func updateStateAfterReadPermission(credentials *apiCustomRole.Permission) PermissionModel {
	var retVal PermissionModel

	retVal.Name = helpers.StringValueOrNull(credentials.Name)

	return retVal
}
