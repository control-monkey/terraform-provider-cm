package stack_dependency

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/stack"
	sdkstack "github.com/control-monkey/controlmonkey-sdk-go/services/stack"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res *sdkstack.Dependency, state *ResourceModel) {
	state.ID = helpers.StringValueOrNull(res.ID)
	state.StackId = helpers.StringValueOrNull(res.StackId)
	state.DependsOnStackId = helpers.StringValueOrNull(res.DependsOnStackId)
	state.TriggerOption = helpers.StringValueOrNull(res.TriggerOption)

	if res.References != nil {
		r := updateStateAfterReadReferences(res.References)
		state.References = r
	} else {
		state.References = nil
	}
}

func updateStateAfterReadReferences(references []*stack.DependencyRef) []*ReferenceModel {
	var retVal []*ReferenceModel

	if references != nil {
		retVal = make([]*ReferenceModel, 0)
		for _, ref := range references {
			r := updateStateAfterReadReference(ref)
			retVal = append(retVal, &r)
		}
	}

	return retVal
}

func updateStateAfterReadReference(ref *stack.DependencyRef) ReferenceModel {
	var retVal ReferenceModel

	retVal.OutputOfStackToDependOn = helpers.StringValueOrNull(ref.OutputOfStackToDependOn)
	retVal.InputForStack = helpers.StringValueOrNull(ref.InputForStack)
	retVal.IncludeSensitiveOutput = helpers.BoolValueOrNull(ref.IncludeSensitiveOutput)

	return retVal
}
