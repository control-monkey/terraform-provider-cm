package template

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	"github.com/control-monkey/controlmonkey-sdk-go/services/template"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*template.Template, bool) {
	var retVal *template.Template

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(template.Template)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if plan.Name != state.Name {
		retVal.SetName(plan.Name.ValueStringPointer())
		hasChanges = true
	}
	if plan.IacType != state.IacType {
		retVal.SetIacType(plan.IacType.ValueStringPointer())
		hasChanges = true
	}
	if plan.Description != state.Description {
		retVal.SetDescription(plan.Description.ValueStringPointer())
		hasChanges = true
	}

	if vcsInfo, hasChanged := vcsInfoConverter(plan.VcsInfo, state.VcsInfo, converterType); hasChanged {
		retVal.SetVcsInfo(vcsInfo)
		hasChanges = true
	}

	if policy, hasChanged := policyConverter(plan.Policy, state.Policy, converterType); hasChanged {
		retVal.SetPolicy(policy)
		hasChanges = true
	}

	if plan.SkipStateRefreshOnDestroy != state.SkipStateRefreshOnDestroy {
		retVal.SetSkipStateRefreshOnDestroy(plan.SkipStateRefreshOnDestroy.ValueBoolPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}

func vcsInfoConverter(plan *VcsInfoModel, state *VcsInfoModel, converterType commons.ConverterType) (*template.VcsInfo, bool) {
	var retVal *template.VcsInfo

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(template.VcsInfo)
	hasChanges := false

	if state == nil {
		state = new(VcsInfoModel) // dummy initialization
		hasChanges = true         // must have changes because before is null and after is not
	}

	if plan.ProviderId != state.ProviderId {
		retVal.SetProviderId(plan.ProviderId.ValueStringPointer())
		hasChanges = true
	}
	if plan.RepoName != state.RepoName {
		retVal.SetRepoName(plan.RepoName.ValueStringPointer())
		hasChanges = true
	}
	if plan.Path != state.Path {
		retVal.SetPath(plan.Path.ValueStringPointer())
		hasChanges = true
	}
	if plan.Branch != state.Branch {
		retVal.SetBranch(plan.Branch.ValueStringPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}

func policyConverter(plan *PolicyModel, state *PolicyModel, converterType commons.ConverterType) (*template.Policy, bool) {
	var retVal *template.Policy

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(template.Policy)
	hasChanges := false

	if state == nil {
		state = new(PolicyModel) // dummy initialization
		hasChanges = true        // must have changes because before is null and after is not
	}

	if innerProperty, hasInnerChanges := ttlConfigConverter(plan.TtlConfig, state.TtlConfig, converterType); hasInnerChanges {
		retVal.SetTtlConfig(innerProperty)
		hasChanges = true
	}
	return retVal, hasChanges
}

func ttlConfigConverter(plan *TtlConfigModel, state *TtlConfigModel, converterType commons.ConverterType) (*template.TtlConfig, bool) {
	var retVal *template.TtlConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(template.TtlConfig)
	hasChanges := false

	if state == nil {
		state = new(TtlConfigModel) // dummy initialization
		hasChanges = true           // must have changes because before is null and after is not
	}

	if innerProperty, hasInnerChanges := ttlDefinitionModelConverter(plan.MaxTtl, state.MaxTtl, converterType); hasInnerChanges {
		retVal.SetMaxTtl(innerProperty)
		hasChanges = true
	}
	if innerProperty, hasInnerChanges := ttlDefinitionModelConverter(plan.DefaultTtl, state.DefaultTtl, converterType); hasInnerChanges {
		retVal.SetDefaultTtl(innerProperty)
		hasChanges = true
	}

	return retVal, hasChanges
}

func ttlDefinitionModelConverter(plan *TtlDefinitionModel, state *TtlDefinitionModel, converterType commons.ConverterType) (*template.TtlDefinition, bool) {
	var retVal *template.TtlDefinition

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(template.TtlDefinition)
	hasChanges := false

	if state == nil {
		state = new(TtlDefinitionModel) // dummy initialization
		hasChanges = true               // must have changes because before is null and after is not
	}

	if plan.Type != state.Type {
		retVal.SetType(plan.Type.ValueStringPointer())
		hasChanges = true
	}
	if plan.Value != state.Value {
		retVal.SetValue(controlmonkey.Int(int(plan.Value.ValueInt64())))
		hasChanges = true
	}

	return retVal, hasChanges
}
