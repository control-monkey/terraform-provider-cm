package stack_dependency

import (
	"reflect"

	sdkStack "github.com/control-monkey/controlmonkey-sdk-go/services/stack"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*sdkStack.Dependency, bool) {
	var retVal *sdkStack.Dependency

	if plan == nil {
		if state == nil {
			return nil, false
		} else {
			return nil, true
		}
	}

	retVal = new(sdkStack.Dependency)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel)
		hasChanges = true
	}

	if plan.StackId != state.StackId {
		retVal.SetStackId(plan.StackId.ValueStringPointer())
		hasChanges = true
	}
	if plan.DependsOnStackId != state.DependsOnStackId {
		retVal.SetDependsOnStackId(plan.DependsOnStackId.ValueStringPointer())
		hasChanges = true
	}
	if plan.TriggerOption != state.TriggerOption {
		retVal.SetTriggerOption(plan.TriggerOption.ValueStringPointer())
		hasChanges = true
	}

	if refs, changed := referencesConverter(plan.References, state.References, converterType); changed {
		retVal.SetReferences(refs)
		hasChanges = true
	}

	return retVal, hasChanges
}

func referencesConverter(plan []*ReferenceModel, state []*ReferenceModel, converterType commons.ConverterType) ([]*sdkStack.DependencyRef, bool) {
	var retVal []*sdkStack.DependencyRef
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true
		if plan != nil {
			retVal = make([]*sdkStack.DependencyRef, 0)
			for _, r := range plan {
				retVal = append(retVal, referenceConverter(r))
			}
		}
	}

	return retVal, hasChanged
}

func referenceConverter(plan *ReferenceModel) *sdkStack.DependencyRef {
	retVal := new(sdkStack.DependencyRef)

	retVal.SetOutputOfStackToDependOn(plan.OutputOfStackToDependOn.ValueStringPointer())
	retVal.SetInputForStack(plan.InputForStack.ValueStringPointer())
	retVal.SetIncludeSensitiveOutput(plan.IncludeSensitiveOutput.ValueBoolPointer())

	return retVal
}
