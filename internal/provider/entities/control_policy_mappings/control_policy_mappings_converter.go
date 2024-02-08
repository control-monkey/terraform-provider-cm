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
	retVal.EntitiesToCreate = convertEntities(mergeResult.EntitiesToCreate, controlPolicyId, createOperation)
	retVal.EntitiesToUpdate = convertEntities(mergeResult.EntitiesToUpdate, controlPolicyId, updateOperation)
	retVal.EntitiesToDelete = convertEntities(mergeResult.EntitiesToDelete, controlPolicyId, deleteOperation)

	return retVal
}

type operationType string

const (
	createOperation operationType = "create"
	updateOperation operationType = "update"
	deleteOperation operationType = "delete"
)

func convertEntities(entities set.Collection[*TargetModel], controlPolicyId types.String, operation operationType) []*controlPolicy.ControlPolicyMapping {
	retVal := make([]*controlPolicy.ControlPolicyMapping, entities.Size())

	for i, e := range entities.Slice() {
		apiEntity := new(controlPolicy.ControlPolicyMapping)
		apiEntity.SetControlPolicyId(controlPolicyId.ValueStringPointer())
		apiEntity.SetTargetId(e.TargetId.ValueStringPointer())
		apiEntity.SetTargetType(e.TargetType.ValueStringPointer())

		if operation != deleteOperation { // enforcement level cannot be sent on delete request
			apiEntity.SetEnforcementLevel(e.EnforcementLevel.ValueStringPointer())
		}

		retVal[i] = apiEntity
	}

	return retVal
}

//
//func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*controlPolicy.ControlPolicyMapping, bool) {
//	var retVal *controlPolicy.ControlPolicyMapping
//
//	if plan == nil {
//		if state == nil {
//			return nil, false // both are the same, no changes
//		} else {
//			return nil, true // before had data, after update is null -> update to null
//		}
//	}
//
//	retVal = new(controlPolicy.ControlPolicyMapping)
//	hasChanges := false
//
//	if state == nil {
//		state = new(ResourceModel) // dummy initialization
//		hasChanges = true          // must have changes because before is null and after is not
//	}
//
//	//The following fields together are the identifier of the resource
//	retVal.SetControlPolicyId(plan.ControlPolicyId.ValueStringPointer())
//	retVal.SetTargetId(plan.TargetId.ValueStringPointer())
//	retVal.SetTargetType(plan.TargetType.ValueStringPointer())
//
//	if plan.EnforcementLevel != state.EnforcementLevel {
//		retVal.SetEnforcementLevel(plan.EnforcementLevel.ValueStringPointer())
//		hasChanges = true
//	}
//
//	return retVal, hasChanges
//}
