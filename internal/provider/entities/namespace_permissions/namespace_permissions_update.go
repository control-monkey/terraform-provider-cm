package namespace_permissions

import (
	namespacePermissions "github.com/control-monkey/controlmonkey-sdk-go/services/namespace_permissions"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(apiEntities []*namespacePermissions.NamespacePermission, state *ResourceModel) {
	state.NamespaceId = state.ID

	if apiEntities != nil {
		permissions := updateStateAfterReadNamespacePermissions(apiEntities)
		state.Permissions = permissions
	} else {
		state.Permissions = nil
	}
}

func updateStateAfterReadNamespacePermissions(permissions []*namespacePermissions.NamespacePermission) []*PermissionsModel {
	var retVal []*PermissionsModel

	if len(permissions) > 0 {
		retVal = make([]*PermissionsModel, len(permissions))

		for i, permission := range permissions {
			u := updateStateAfterReadPermission(permission)
			retVal[i] = &u
		}
	}

	return retVal
}

func updateStateAfterReadPermission(permission *namespacePermissions.NamespacePermission) PermissionsModel {
	var retVal PermissionsModel

	retVal.UserEmail = helpers.StringValueOrNull(permission.UserEmail)
	retVal.ProgrammaticUserName = helpers.StringValueOrNull(permission.ProgrammaticUserName)
	retVal.TeamId = helpers.StringValueOrNull(permission.TeamId)
	retVal.Role = helpers.StringValueOrNull(permission.Role)
	retVal.CustomRoleId = helpers.StringValueOrNull(permission.CustomRoleId)

	return retVal
}
