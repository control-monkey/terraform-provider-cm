package customRole

import (
	apiCustomRole "github.com/control-monkey/controlmonkey-sdk-go/services/custom_role"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(res *apiCustomRole.CustomRole, state *ResourceModel) {
	state.Type = helpers.StringValueOrNull(res.Type)
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

	if credentials.Names != nil {
		retVal.Names = make([]types.String, 0)
		for _, n := range credentials.Names {
			retVal.Names = append(retVal.Names, helpers.StringValueOrNull(n))
		}
	}

	if credentials.Restrictions != nil {
		retVal.Restrictions = make([]*RestrictionModel, 0)
		for _, r := range credentials.Restrictions {
			restriction := updateStateAfterReadRestriction(r)
			retVal.Restrictions = append(retVal.Restrictions, &restriction)
		}
	}

	return retVal
}

func updateStateAfterReadRestriction(restriction *apiCustomRole.PermissionRestriction) RestrictionModel {
	var retVal RestrictionModel

	retVal.CloudProvider = helpers.StringValueOrNull(restriction.CloudProvider)
	retVal.CloudAccountId = helpers.StringValueOrNull(restriction.CloudAccountId)
	retVal.CmResourceName = helpers.StringValueOrNull(restriction.CmResourceName)

	return retVal
}
