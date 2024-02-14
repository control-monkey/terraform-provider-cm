package controlPolicyGroupMappings

import (
	controlPolicyGroup "github.com/control-monkey/controlmonkey-sdk-go/services/control_policy_group"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	"github.com/hashicorp/go-set/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MergedEntities struct {
	EntitiesToCreate []*controlPolicyGroup.ControlPolicyGroupMapping
	EntitiesToUpdate []*controlPolicyGroup.ControlPolicyGroupMapping
	EntitiesToDelete []*controlPolicyGroup.ControlPolicyGroupMapping
}

func Merge(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) *MergedEntities {
	retVal := new(MergedEntities)

	if plan == nil {
		plan = new(ResourceModel) // delete merger
	}

	if state == nil {
		state = new(ResourceModel) // create merger
	}

	var controlPolicyGroupId types.String
	if plan.ControlPolicyGroupId.IsNull() == false {
		controlPolicyGroupId = plan.ControlPolicyGroupId
	} else {
		controlPolicyGroupId = state.ControlPolicyGroupId
	}

	mergeResult := interfaces.MergeEntities(plan.Targets, state.Targets)
	retVal.EntitiesToCreate = convertEntities(mergeResult.EntitiesToCreate, controlPolicyGroupId, interfaces.CreateOperation)
	retVal.EntitiesToUpdate = convertEntities(mergeResult.EntitiesToUpdate, controlPolicyGroupId, interfaces.UpdateOperation)
	retVal.EntitiesToDelete = convertEntities(mergeResult.EntitiesToDelete, controlPolicyGroupId, interfaces.DeleteOperation)

	return retVal
}

func convertEntities(entities set.Collection[*TargetModel], controlPolicyGroupId types.String, operation interfaces.OperationType) []*controlPolicyGroup.ControlPolicyGroupMapping {
	retVal := make([]*controlPolicyGroup.ControlPolicyGroupMapping, entities.Size())

	for i, e := range entities.Slice() {
		apiEntity := new(controlPolicyGroup.ControlPolicyGroupMapping)
		apiEntity.SetControlPolicyGroupId(controlPolicyGroupId.ValueStringPointer())
		apiEntity.SetTargetId(e.TargetId.ValueStringPointer())
		apiEntity.SetTargetType(e.TargetType.ValueStringPointer())

		if operation != interfaces.DeleteOperation { // enforcement level cannot be sent on delete request
			apiEntity.SetEnforcementLevel(e.EnforcementLevel.ValueStringPointer())

			apiOverrideEnforcements := convertOverrideEnforcements(e.OverrideEnforcements)
			apiEntity.SetOverrideEnforcements(apiOverrideEnforcements)
		}

		retVal[i] = apiEntity
	}

	return retVal
}

func convertOverrideEnforcements(es []*OverrideEnforcementModel) []*controlPolicyGroup.OverrideEnforcement {
	retVal := make([]*controlPolicyGroup.OverrideEnforcement, len(es))

	for i, e := range es {
		apiEntity := convertOverrideEnforcement(e)
		retVal[i] = apiEntity
	}

	return retVal
}

func convertOverrideEnforcement(e *OverrideEnforcementModel) *controlPolicyGroup.OverrideEnforcement {
	retVal := new(controlPolicyGroup.OverrideEnforcement)

	retVal.SetControlPolicyId(e.ControlPolicyId.ValueStringPointer())
	retVal.SetEnforcementLevel(e.EnforcementLevel.ValueStringPointer())

	apiStackIds := helpers.TfListToStringSlice(e.StackIds)
	retVal.SetStackIds(apiStackIds)

	return retVal
}
