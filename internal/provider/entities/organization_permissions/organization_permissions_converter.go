package organizationPermissions

import (
	organizationPermissions "github.com/control-monkey/controlmonkey-sdk-go/services/organization_permissions"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	"github.com/hashicorp/go-set/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MergedEntities struct {
	EntitiesToCreate []*organizationPermissions.OrganizationPermission
	EntitiesToDelete []*organizationPermissions.OrganizationPermission
}

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

	mergeResult := interfaces.MergeEntities(plan.Permissions, state.Permissions)
	retVal.EntitiesToCreate = convertEntities(mergeResult.EntitiesToCreate, teamId)
	retVal.EntitiesToDelete = convertEntities(mergeResult.EntitiesToDelete, teamId)

	return retVal
}

func convertEntities(entities set.Collection[*PermissionModel], teamId types.String) []*organizationPermissions.OrganizationPermission {
	retVal := make([]*organizationPermissions.OrganizationPermission, entities.Size())

	for i, e := range entities.Slice() {
		apiEntity := new(organizationPermissions.OrganizationPermission)
		apiEntity.SetTeamId(teamId.ValueStringPointer())
		apiEntity.SetCustomRoleId(e.CustomRoleId.ValueStringPointer())

		retVal[i] = apiEntity
	}

	return retVal
}
