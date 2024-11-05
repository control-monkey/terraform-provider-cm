package customAbacConfiguration

import (
	apiCustomAbacConfiguration "github.com/control-monkey/controlmonkey-sdk-go/services/custom_abac_configuration"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"reflect"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*apiCustomAbacConfiguration.CustomAbacConfiguration, bool) {
	var retVal *apiCustomAbacConfiguration.CustomAbacConfiguration

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(apiCustomAbacConfiguration.CustomAbacConfiguration)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if plan.Name != state.Name {
		retVal.SetName(plan.Name.ValueStringPointer())
		hasChanges = true
	}
	if plan.CustomAbacId != state.CustomAbacId {
		retVal.SetCustomAbacId(plan.CustomAbacId.ValueStringPointer())
		hasChanges = true
	}
	if cp, hasChanged := RolesConverter(plan.Roles, state.Roles, converterType); hasChanged {
		retVal.SetRoles(cp)
		hasChanges = true
	}

	return retVal, hasChanges
}

func RolesConverter(plan []*RoleModel, state []*RoleModel, converterType commons.ConverterType) ([]*apiCustomAbacConfiguration.Role, bool) {
	var retVal []*apiCustomAbacConfiguration.Role
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true

		if plan != nil {
			retVal = make([]*apiCustomAbacConfiguration.Role, 0)

			for _, r := range plan {
				role := roleConverter(r)
				retVal = append(retVal, role)
			}
		}
	}

	return retVal, hasChanged
}

func roleConverter(plan *RoleModel) *apiCustomAbacConfiguration.Role {
	retVal := new(apiCustomAbacConfiguration.Role)

	retVal.SetOrgId(plan.OrgId.ValueStringPointer())
	retVal.SetOrgRole(plan.OrgRole.ValueStringPointer())
	retVal.SetTeamIds(helpers.TfListToStringSlice(plan.TeamIds))

	return retVal
}
