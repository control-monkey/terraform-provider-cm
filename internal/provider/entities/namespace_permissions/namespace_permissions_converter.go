package namespace_permissions

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/namespace_permissions"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/hashicorp/go-set/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

	planEntities := set.HashSetFrom[*PermissionsModel, string](plan.Permissions)
	stateEntities := set.HashSetFrom[*PermissionsModel, string](state.Permissions)

	permissionsToCreate := planEntities.Difference(stateEntities)
	permissionsToDelete := stateEntities.Difference(planEntities)

	idToModelToDelete := make(map[string]*PermissionsModel, 0)
	permissionsToRemoveFromCreate := make([]*PermissionsModel, 0)
	permissionsToRemoveFromDelete := make([]*PermissionsModel, 0)
	permissionsToUpdate := set.NewHashSet[*PermissionsModel, string](0)

	for _, p := range permissionsToDelete.Slice() {
		idToModelToDelete[p.GetBlockIdentifier()] = p
	}
	for _, p := range permissionsToCreate.Slice() {
		if permissionToDelete := idToModelToDelete[p.GetBlockIdentifier()]; permissionToDelete != nil { // id in both add & delete
			permissionsToRemoveFromCreate = append(permissionsToRemoveFromCreate, p)
			permissionsToRemoveFromDelete = append(permissionsToRemoveFromDelete, permissionToDelete)
			permissionsToUpdate.Insert(p)
		}
	}

	permissionsToCreate.RemoveSlice(permissionsToRemoveFromCreate)
	permissionsToDelete.RemoveSlice(permissionsToRemoveFromDelete)

	retVal.EntitiesToCreate = buildEntities(permissionsToCreate, namespaceId)
	retVal.EntitiesToUpdate = buildEntities(permissionsToUpdate, namespaceId)
	retVal.EntitiesToDelete = buildEntities(permissionsToDelete, namespaceId)

	return retVal
}

func buildEntities(entities set.Collection[*PermissionsModel], namespaceId types.String) []*namespace_permissions.NamespacePermission {
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
