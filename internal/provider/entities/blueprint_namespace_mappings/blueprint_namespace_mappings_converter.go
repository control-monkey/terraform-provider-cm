package blueprint_namespace_mappings

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/blueprint"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	"github.com/hashicorp/go-set/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MergedEntities struct {
	EntitiesToCreate []*blueprint.BlueprintNamespaceMapping
	EntitiesToUpdate []*blueprint.BlueprintNamespaceMapping
	EntitiesToDelete []*blueprint.BlueprintNamespaceMapping
}

func Merge(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) *MergedEntities {
	retVal := new(MergedEntities)

	if plan == nil {
		plan = new(ResourceModel) // delete merger
	}

	if state == nil {
		state = new(ResourceModel) // create merger
	}

	var blueprintId types.String
	if plan.BlueprintId.IsNull() == false {
		blueprintId = plan.BlueprintId
	} else {
		blueprintId = state.BlueprintId
	}

	mergeResult := interfaces.MergeEntities(plan.Namespaces, state.Namespaces)
	retVal.EntitiesToCreate = convertEntities(mergeResult.EntitiesToCreate, blueprintId)
	retVal.EntitiesToUpdate = convertEntities(mergeResult.EntitiesToUpdate, blueprintId)
	retVal.EntitiesToDelete = convertEntities(mergeResult.EntitiesToDelete, blueprintId)

	return retVal
}

func convertEntities(entities set.Collection[*NamespaceModel], blueprintId types.String) []*blueprint.BlueprintNamespaceMapping {
	retVal := make([]*blueprint.BlueprintNamespaceMapping, entities.Size())

	for i, u := range entities.Slice() {
		tu := new(blueprint.BlueprintNamespaceMapping)
		tu.SetBlueprintId(blueprintId.ValueStringPointer())
		tu.SetNamespaceId(u.NamespaceId.ValueStringPointer())

		retVal[i] = tu
	}

	return retVal
}
