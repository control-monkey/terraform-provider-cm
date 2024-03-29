package controlPolicyMappings

import (
	controlPolicy "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	"github.com/hashicorp/go-set/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MergedEntities struct {
	EntitiesToCreate []*controlPolicy.ControlPolicyMapping
	EntitiesToUpdate []*controlPolicy.ControlPolicyMapping
	EntitiesToDelete []*controlPolicy.ControlPolicyMapping
}

func Merge(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) *MergedEntities {
	retVal := new(MergedEntities)

	if plan == nil {
		plan = new(ResourceModel) // delete merger
	}

	if state == nil {
		state = new(ResourceModel) // create merger
	}

	var controlPolicyId types.String
	if plan.ControlPolicyId.IsNull() == false {
		controlPolicyId = plan.ControlPolicyId
	} else {
		controlPolicyId = state.ControlPolicyId
	}

	mergeResult := interfaces.MergeEntities(plan.Targets, state.Targets)
	retVal.EntitiesToCreate = convertEntities(mergeResult.EntitiesToCreate, controlPolicyId, interfaces.CreateOperation)
	retVal.EntitiesToUpdate = convertEntities(mergeResult.EntitiesToUpdate, controlPolicyId, interfaces.UpdateOperation)
	retVal.EntitiesToDelete = convertEntities(mergeResult.EntitiesToDelete, controlPolicyId, interfaces.DeleteOperation)

	return retVal
}

func convertEntities(entities set.Collection[*TargetModel], controlPolicyId types.String, operation interfaces.OperationType) []*controlPolicy.ControlPolicyMapping {
	retVal := make([]*controlPolicy.ControlPolicyMapping, entities.Size())

	for i, e := range entities.Slice() {
		apiEntity := new(controlPolicy.ControlPolicyMapping)
		apiEntity.SetControlPolicyId(controlPolicyId.ValueStringPointer())
		apiEntity.SetTargetId(e.TargetId.ValueStringPointer())
		apiEntity.SetTargetType(e.TargetType.ValueStringPointer())

		if operation != interfaces.DeleteOperation { // enforcement level cannot be sent on delete request
			apiEntity.SetEnforcementLevel(e.EnforcementLevel.ValueStringPointer())
		}

		retVal[i] = apiEntity
	}

	return retVal
}
