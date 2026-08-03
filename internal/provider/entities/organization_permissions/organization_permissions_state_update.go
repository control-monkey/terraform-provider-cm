package organizationPermissions

import (
	sdkOrganizationPermissions "github.com/control-monkey/controlmonkey-sdk-go/services/organization_permissions"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(apiEntities []*sdkOrganizationPermissions.OrganizationPermission, state *ResourceModel) {
	state.TeamId = state.ID

	if apiEntities != nil {
		permissions := updateStateAfterReadPermissions(apiEntities)
		state.Permissions = permissions
	} else {
		state.Permissions = nil
	}
}

func updateStateAfterReadPermissions(permissions []*sdkOrganizationPermissions.OrganizationPermission) []*PermissionModel {
	var retVal []*PermissionModel

	if len(permissions) > 0 {
		retVal = make([]*PermissionModel, len(permissions))

		for i, permission := range permissions {
			u := updateStateAfterReadPermission(permission)
			retVal[i] = &u
		}
	}

	return retVal
}

func updateStateAfterReadPermission(permission *sdkOrganizationPermissions.OrganizationPermission) PermissionModel {
	var retVal PermissionModel

	retVal.CustomRoleId = helpers.StringValueOrNull(permission.CustomRoleId)

	return retVal
}
