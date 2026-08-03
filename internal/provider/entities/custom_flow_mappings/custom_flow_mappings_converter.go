package customFlowMappings

import (
	customFlow "github.com/control-monkey/controlmonkey-sdk-go/services/custom_flow"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	"github.com/hashicorp/go-set/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MergedEntities struct {
	EntitiesToCreate []*customFlow.CustomFlowMapping
	EntitiesToDelete []*customFlow.CustomFlowMapping
}

func Merge(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) *MergedEntities {
	retVal := new(MergedEntities)

	if plan == nil {
		plan = new(ResourceModel) // delete merger
	}

	if state == nil {
		state = new(ResourceModel) // create merger
	}

	var customFlowId types.String
	if plan.CustomFlowId.IsNull() == false {
		customFlowId = plan.CustomFlowId
	} else {
		customFlowId = state.CustomFlowId
	}

	mergeResult := interfaces.MergeEntities(plan.Targets, state.Targets)
	retVal.EntitiesToCreate = convertEntities(mergeResult.EntitiesToCreate, customFlowId)
	retVal.EntitiesToDelete = convertEntities(mergeResult.EntitiesToDelete, customFlowId)

	return retVal
}

func convertEntities(entities set.Collection[*TargetModel], customFlowId types.String) []*customFlow.CustomFlowMapping {
	retVal := make([]*customFlow.CustomFlowMapping, entities.Size())

	for i, e := range entities.Slice() {
		apiEntity := new(customFlow.CustomFlowMapping)
		apiEntity.SetCustomFlowId(customFlowId.ValueStringPointer())
		apiEntity.SetTargetId(e.TargetId.ValueStringPointer())
		apiEntity.SetTargetType(e.TargetType.ValueStringPointer())

		retVal[i] = apiEntity
	}

	return retVal
}
