package blueprint

import (
	"reflect"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	apiBlueprint "github.com/control-monkey/controlmonkey-sdk-go/services/blueprint"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*apiBlueprint.Blueprint, bool) {
	var retVal *apiBlueprint.Blueprint

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(apiBlueprint.Blueprint)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if plan.Name != state.Name {
		retVal.SetName(plan.Name.ValueStringPointer())
		hasChanges = true
	}

	if plan.Description != state.Description {
		retVal.SetDescription(plan.Description.ValueStringPointer())
		hasChanges = true
	}

	if vcsInfo, hasChanged := blueprintVcsInfoConverter(plan.BlueprintVcsInfo, state.BlueprintVcsInfo, converterType); hasChanged {
		retVal.SetBlueprintVcsInfo(vcsInfo)
		hasChanges = true
	}

	if sc, hasChanged := stackConfigurationConverter(plan.StackConfiguration, state.StackConfiguration, converterType); hasChanged {
		retVal.SetStackConfiguration(sc)
		hasChanges = true
	}

	if sp, hasChanged := substituteParametersConverter(plan.SubstituteParameters, state.SubstituteParameters, converterType); hasChanged {
		retVal.SetSubstituteParameters(sp)
		hasChanges = true
	}

	if plan.SkipPlanOnStackInitialization != state.SkipPlanOnStackInitialization {
		retVal.SetSkipPlanOnStackInitialization(plan.SkipPlanOnStackInitialization.ValueBoolPointer())
		hasChanges = true
	}
	if plan.AutoApproveApplyOnInitialization != state.AutoApproveApplyOnInitialization {
		retVal.SetAutoApproveApplyOnInitialization(plan.AutoApproveApplyOnInitialization.ValueBoolPointer())
		hasChanges = true
	}

	if policy, hasChanged := policyConverter(plan.Policy, state.Policy, converterType); hasChanged {
		retVal.SetPolicy(policy)
		hasChanges = true
	}

	return retVal, hasChanges
}

func blueprintVcsInfoConverter(plan *VcsInfoModel, state *VcsInfoModel, converterType commons.ConverterType) (*apiBlueprint.VcsInfo, bool) {
	var retVal *apiBlueprint.VcsInfo

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(apiBlueprint.VcsInfo)
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

func stackConfigurationConverter(plan *StackConfigurationModel, state *StackConfigurationModel, converterType commons.ConverterType) (*apiBlueprint.StackConfiguration, bool) {
	var retVal *apiBlueprint.StackConfiguration

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(apiBlueprint.StackConfiguration)
	hasChanges := false

	if state == nil {
		state = new(StackConfigurationModel) // dummy initialization
		hasChanges = true                    // must have changes because before is null and after is not
	}

	if plan.NamePattern != state.NamePattern {
		retVal.SetNamePattern(plan.NamePattern.ValueStringPointer())
		hasChanges = true
	}
	if plan.IacType != state.IacType {
		retVal.SetIacType(plan.IacType.ValueStringPointer())
		hasChanges = true
	}
	if vcsInfo, hasChanged := vcsInfoWithPatternsConverter(plan.VcsInfoWithPatterns, state.VcsInfoWithPatterns, converterType); hasChanged {
		retVal.SetVcsInfoWithPatterns(vcsInfo)
		hasChanges = true
	}

	if deploymentApprovalPolicy, hasChanged := cross_models.DeploymentApprovalPolicyConverter(plan.DeploymentApprovalPolicy, state.DeploymentApprovalPolicy, converterType); hasChanged {
		retVal.SetDeploymentApprovalPolicy(deploymentApprovalPolicy)
		hasChanges = true
	}

	if runTrigger, hasChanged := cross_models.RunTriggerConverter(plan.RunTrigger, state.RunTrigger, converterType); hasChanged {
		retVal.SetRunTrigger(runTrigger)
		hasChanges = true
	}

	if iacConfig, hasChanged := cross_models.IacConfigConverter(plan.IacConfig, state.IacConfig, converterType); hasChanged {
		retVal.SetIacConfig(iacConfig)
		hasChanges = true
	}

	if autoSync, hasChanged := cross_models.AutoSyncConverter(plan.AutoSync, state.AutoSync, converterType); hasChanged {
		retVal.SetAutoSync(autoSync)
		hasChanges = true
	}

	return retVal, hasChanges
}

func vcsInfoWithPatternsConverter(plan *StackVcsInfoWithPatternsModel, state *StackVcsInfoWithPatternsModel, converterType commons.ConverterType) (*apiBlueprint.StackVcsInfoWithPatterns, bool) {
	var retVal *apiBlueprint.StackVcsInfoWithPatterns

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(apiBlueprint.StackVcsInfoWithPatterns)
	hasChanges := false

	if state == nil {
		state = new(StackVcsInfoWithPatternsModel) // dummy initialization
		hasChanges = true                          // must have changes because before is null and after is not
	}

	if plan.ProviderId != state.ProviderId {
		retVal.SetProviderId(plan.ProviderId.ValueStringPointer())
		hasChanges = true
	}
	if plan.RepoName != state.RepoName {
		retVal.SetRepoName(plan.RepoName.ValueStringPointer())
		hasChanges = true
	}
	if plan.PathPattern != state.PathPattern {
		retVal.SetPathPattern(plan.PathPattern.ValueStringPointer())
		hasChanges = true
	}
	if plan.BranchPattern != state.BranchPattern {
		retVal.SetBranchPattern(plan.BranchPattern.ValueStringPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}

func substituteParametersConverter(plan []*SubstituteParameterModel, state []*SubstituteParameterModel, converterType commons.ConverterType) ([]*apiBlueprint.SubstituteParameter, bool) {
	var retVal []*apiBlueprint.SubstituteParameter
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true

		if plan != nil {
			retVal = make([]*apiBlueprint.SubstituteParameter, 0)

			for _, r := range plan {
				substituteParameter := substituteParameterConverter(r, converterType)
				retVal = append(retVal, substituteParameter)
			}
		}
	}

	return retVal, hasChanged
}

func substituteParameterConverter(plan *SubstituteParameterModel, converterType commons.ConverterType) *apiBlueprint.SubstituteParameter {
	retVal := new(apiBlueprint.SubstituteParameter)

	retVal.SetKey(plan.Key.ValueStringPointer())
	retVal.SetDescription(plan.Description.ValueStringPointer())

	vc, _ := cross_models.ValueConditionsConverter(plan.ValueConditions, nil, converterType)
	retVal.SetValueConditions(vc)

	return retVal
}

func policyConverter(plan *PolicyModel, state *PolicyModel, converterType commons.ConverterType) (*apiBlueprint.Policy, bool) {
	var retVal *apiBlueprint.Policy

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(apiBlueprint.Policy)
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

func ttlConfigConverter(plan *TtlConfigModel, state *TtlConfigModel, converterType commons.ConverterType) (*apiBlueprint.TtlConfig, bool) {
	var retVal *apiBlueprint.TtlConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(apiBlueprint.TtlConfig)
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
	if plan.OpenCleanupPrOnTtlTermination != state.OpenCleanupPrOnTtlTermination {
		retVal.SetOpenCleanupPrOnTtlTermination(plan.OpenCleanupPrOnTtlTermination.ValueBoolPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}

func ttlDefinitionModelConverter(plan *TtlDefinitionModel, state *TtlDefinitionModel, converterType commons.ConverterType) (*apiBlueprint.TtlDefinition, bool) {
	var retVal *apiBlueprint.TtlDefinition

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(apiBlueprint.TtlDefinition)
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
