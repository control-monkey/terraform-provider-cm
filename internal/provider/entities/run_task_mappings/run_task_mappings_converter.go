package runTaskMappings

import (
	runTask "github.com/control-monkey/controlmonkey-sdk-go/services/run_task"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	"github.com/hashicorp/go-set/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MergedEntities struct {
	EntitiesToCreate []*runTask.RunTaskMapping
	EntitiesToUpdate []*runTask.RunTaskMapping
	EntitiesToDelete []*runTask.RunTaskMapping
}

func Merge(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) *MergedEntities {
	retVal := new(MergedEntities)

	if plan == nil {
		plan = new(ResourceModel) // delete merger
	}

	if state == nil {
		state = new(ResourceModel) // create merger
	}

	var runTaskId types.String
	if plan.RunTaskId.IsNull() == false {
		runTaskId = plan.RunTaskId
	} else {
		runTaskId = state.RunTaskId
	}

	mergeResult := interfaces.MergeEntities(plan.Targets, state.Targets)
	retVal.EntitiesToCreate = convertEntities(mergeResult.EntitiesToCreate, runTaskId, interfaces.CreateOperation)
	retVal.EntitiesToUpdate = convertEntities(mergeResult.EntitiesToUpdate, runTaskId, interfaces.UpdateOperation)
	retVal.EntitiesToDelete = convertEntities(mergeResult.EntitiesToDelete, runTaskId, interfaces.DeleteOperation)

	return retVal
}

func convertEntities(entities set.Collection[*TargetModel], runTaskId types.String, operation interfaces.OperationType) []*runTask.RunTaskMapping {
	retVal := make([]*runTask.RunTaskMapping, entities.Size())

	for i, e := range entities.Slice() {
		apiEntity := new(runTask.RunTaskMapping)
		apiEntity.SetRunTaskId(runTaskId.ValueStringPointer())
		apiEntity.SetTargetId(e.TargetId.ValueStringPointer())
		apiEntity.SetTargetType(e.TargetType.ValueStringPointer())

		if operation != interfaces.DeleteOperation { // enforcement level and stage cannot be sent on delete request
			apiEntity.SetEnforcementLevel(e.EnforcementLevel.ValueStringPointer())
			apiEntity.SetStage(e.Stage.ValueStringPointer())
		}

		retVal[i] = apiEntity
	}

	return retVal
}
