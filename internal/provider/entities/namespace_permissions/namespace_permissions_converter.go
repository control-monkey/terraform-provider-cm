package namespace_permissions

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/namespace_permissions"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	"github.com/hashicorp/go-set/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MergedEntities struct {
	EntitiesToCreate []*namespace_permissions.NamespacePermission
	EntitiesToUpdate []*namespace_permissions.NamespacePermission
	EntitiesToDelete []*namespace_permissions.NamespacePermission
}

func Merge(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) *MergedEntities {
	retVal := new(MergedEntities)

	if plan == nil {
		plan = new(ResourceModel) // delete merger
	}

	if state == nil {
		state = new(ResourceModel) // create merger
	}

	var namespaceId types.String
	if plan.NamespaceId.IsNull() == false {
		namespaceId = plan.NamespaceId
	} else {
		namespaceId = state.NamespaceId
	}

	mergeResult := interfaces.MergeEntities(plan.Permissions, state.Permissions)
	retVal.EntitiesToCreate = convertEntities(mergeResult.EntitiesToCreate, namespaceId)
	retVal.EntitiesToUpdate = convertEntities(mergeResult.EntitiesToUpdate, namespaceId)
	retVal.EntitiesToDelete = convertEntities(mergeResult.EntitiesToDelete, namespaceId)

	return retVal
}

func convertEntities(entities set.Collection[*PermissionsModel], namespaceId types.String) []*namespace_permissions.NamespacePermission {
	retVal := make([]*namespace_permissions.NamespacePermission, entities.Size())

	for i, e := range entities.Slice() {
		apiEntity := new(namespace_permissions.NamespacePermission)
		apiEntity.SetNamespaceId(namespaceId.ValueStringPointer())
		apiEntity.SetUserEmail(e.UserEmail.ValueStringPointer())
		apiEntity.SetProgrammaticUserName(e.ProgrammaticUserName.ValueStringPointer())
		apiEntity.SetTeamId(e.TeamId.ValueStringPointer())
		apiEntity.SetRole(e.Role.ValueStringPointer())
		apiEntity.SetCustomRoleId(e.CustomRoleId.ValueStringPointer())

		retVal[i] = apiEntity
	}

	return retVal
}
