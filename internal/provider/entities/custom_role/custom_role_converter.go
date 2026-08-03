package customRole

import (
	"reflect"

	apiCustomRole "github.com/control-monkey/controlmonkey-sdk-go/services/custom_role"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*apiCustomRole.CustomRole, bool) {
	var retVal *apiCustomRole.CustomRole

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(apiCustomRole.CustomRole)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if plan.Type != state.Type {
		retVal.SetType(plan.Type.ValueStringPointer())
		hasChanges = true
	}
	if plan.Name != state.Name {
		retVal.SetName(plan.Name.ValueStringPointer())
		hasChanges = true
	}
	if plan.Description != state.Description {
		retVal.SetDescription(plan.Description.ValueStringPointer())
		hasChanges = true
	}
	if cp, hasChanged := PermissionsConverter(plan.Permissions, state.Permissions, converterType); hasChanged {
		retVal.SetPermissions(cp)
		hasChanges = true
	}
	if plan.StackRestriction != state.StackRestriction {
		retVal.SetStackRestriction(plan.StackRestriction.ValueStringPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}

func PermissionsConverter(plan []*PermissionModel, state []*PermissionModel, converterType commons.ConverterType) ([]*apiCustomRole.Permission, bool) {
	var retVal []*apiCustomRole.Permission
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true

		if plan != nil {
			retVal = make([]*apiCustomRole.Permission, 0)

			for _, r := range plan {
				permission := permissionConverter(r)
				retVal = append(retVal, permission)
			}
		}
	}

	return retVal, hasChanged
}

func permissionConverter(plan *PermissionModel) *apiCustomRole.Permission {
	retVal := new(apiCustomRole.Permission)

	// name and names are mutually exclusive, so only set what was configured.
	if plan.Name.IsNull() == false {
		retVal.SetName(plan.Name.ValueStringPointer())
	}
	if plan.Names != nil {
		names := make([]*string, 0)
		for _, n := range plan.Names {
			names = append(names, n.ValueStringPointer())
		}
		retVal.SetNames(names)
	}
	if plan.Restrictions != nil {
		restrictions := make([]*apiCustomRole.PermissionRestriction, 0)
		for _, r := range plan.Restrictions {
			restrictions = append(restrictions, restrictionConverter(r))
		}
		retVal.SetRestrictions(restrictions)
	}

	return retVal
}

func restrictionConverter(plan *RestrictionModel) *apiCustomRole.PermissionRestriction {
	retVal := new(apiCustomRole.PermissionRestriction)

	retVal.SetCloudProvider(plan.CloudProvider.ValueStringPointer())
	retVal.SetCloudAccountId(plan.CloudAccountId.ValueStringPointer())
	retVal.SetCmResourceName(plan.CmResourceName.ValueStringPointer())

	return retVal
}
