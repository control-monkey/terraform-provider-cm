package stack

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkStack "github.com/control-monkey/controlmonkey-sdk-go/services/stack"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*sdkStack.Stack, bool) {
	var retVal *sdkStack.Stack

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(sdkStack.Stack)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if plan.NamespaceId != state.NamespaceId {
		retVal.SetNamespaceId(plan.NamespaceId.ValueStringPointer())
		hasChanges = true
	}
	if plan.IacType != state.IacType {
		retVal.SetIacType(plan.IacType.ValueStringPointer())
		hasChanges = true
	}
	if plan.Name != state.Name {
		retVal.SetName(plan.Name.ValueStringPointer())
		hasChanges = true
	}
	if plan.Description != state.Description {
		retVal.SetDescription(plan.Description.ValueStringPointer())
		hasChanges = true
	}

	var data sdkStack.Data

	if deploymentBehavior, hasChanged := cross_models.DeploymentBehaviorConverter(plan.DeploymentBehavior, state.DeploymentBehavior, converterType); hasChanged {
		data.SetDeploymentBehavior(deploymentBehavior)
		hasChanges = true
	}

	if deploymentApprovalPolicy, hasChanged := cross_models.DeploymentApprovalPolicyConverter(plan.DeploymentApprovalPolicy, state.DeploymentApprovalPolicy, converterType); hasChanged {
		data.SetDeploymentApprovalPolicy(deploymentApprovalPolicy)
		hasChanges = true
	}

	if vcsInfo, hasChanged := vcsInfoConverter(plan.VcsInfo, state.VcsInfo, converterType); hasChanged {
		data.SetVcsInfo(vcsInfo)
		hasChanges = true
	}

	if runTrigger, hasChanged := cross_models.RunTriggerConverter(plan.RunTrigger, state.RunTrigger, converterType); hasChanged {
		data.SetRunTrigger(runTrigger)
		hasChanges = true
	}

	if iacConfig, hasChanged := cross_models.IacConfigConverter(plan.IacConfig, state.IacConfig, converterType); hasChanged {
		data.SetIacConfig(iacConfig)
		hasChanges = true
	}

	if policy, hasChanged := policyConverter(plan.Policy, state.Policy, converterType); hasChanged {
		data.SetPolicy(policy)
		hasChanges = true
	}

	if runnerConfig, hasChanged := cross_models.RunnerConfigConverter(plan.RunnerConfig, state.RunnerConfig, converterType); hasChanged {
		data.SetRunnerConfig(runnerConfig)
		hasChanges = true
	}

	if autoSync, hasChanged := cross_models.AutoSyncConverter(plan.AutoSync, state.AutoSync, converterType); hasChanged {
		data.SetAutoSync(autoSync)
		hasChanges = true
	}

	if capabilities, hasChanged := capabilitiesConverter(plan.Capabilities, state.Capabilities, converterType); hasChanged {
		data.SetCapabilities(capabilities)
		hasChanges = true
	}

	retVal.SetData(&data)

	return retVal, hasChanges
}

func vcsInfoConverter(plan *VcsInfoModel, state *VcsInfoModel, converterType commons.ConverterType) (*sdkStack.VcsInfo, bool) {
	var retVal *sdkStack.VcsInfo

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(sdkStack.VcsInfo)
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

func policyConverter(plan *PolicyModel, state *PolicyModel, converterType commons.ConverterType) (*sdkStack.Policy, bool) {
	var retVal *sdkStack.Policy

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(sdkStack.Policy)
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

func ttlConfigConverter(plan *TtlConfigModel, state *TtlConfigModel, converterType commons.ConverterType) (*sdkStack.TtlConfig, bool) {
	var retVal *sdkStack.TtlConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(sdkStack.TtlConfig)
	hasChanges := false

	if state == nil {
		state = new(TtlConfigModel) // dummy initialization
		hasChanges = true           // must have changes because before is null and after is not
	}

	if innerProperty, hasInnerChanges := ttlDefinitionModelConverter(plan.Ttl, state.Ttl, converterType); hasInnerChanges {
		retVal.SetTtl(innerProperty)
		hasChanges = true
	}
	return retVal, hasChanges
}

func ttlDefinitionModelConverter(plan *TtlDefinitionModel, state *TtlDefinitionModel, converterType commons.ConverterType) (*sdkStack.TtlDefinition, bool) {
	var retVal *sdkStack.TtlDefinition

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(sdkStack.TtlDefinition)
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

func capabilitiesConverter(plan *CapabilitiesModel, state *CapabilitiesModel, converterType commons.ConverterType) (*sdkStack.Capabilities, bool) {
	var retVal *sdkStack.Capabilities

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(sdkStack.Capabilities)
	hasChanges := false

	if state == nil {
		state = new(CapabilitiesModel) // dummy initialization
		hasChanges = true              // must have changes because before is null and after is not
	}

	if deployOnPush, hasChanged := capabilityConfigConverter(plan.DeployOnPush, state.DeployOnPush, converterType); hasChanged {
		retVal.SetDeployOnPush(deployOnPush)
		hasChanges = true
	}

	if planOnPr, hasChanged := capabilityConfigConverter(plan.PlanOnPr, state.PlanOnPr, converterType); hasChanged {
		retVal.SetPlanOnPr(planOnPr)
		hasChanges = true
	}

	if driftDetection, hasChanged := capabilityConfigConverter(plan.DriftDetection, state.DriftDetection, converterType); hasChanged {
		retVal.SetDriftDetection(driftDetection)
		hasChanges = true
	}

	return retVal, hasChanges
}

func capabilityConfigConverter(plan *CapabilityConfigModel, state *CapabilityConfigModel, converterType commons.ConverterType) (*sdkStack.CapabilityConfig, bool) {
	var retVal *sdkStack.CapabilityConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(sdkStack.CapabilityConfig)
	hasChanges := false

	if state == nil {
		state = new(CapabilityConfigModel) // dummy initialization
		hasChanges = true                  // must have changes because before is null and after is not
	}

	if plan.Status != state.Status {
		retVal.SetStatus(plan.Status.ValueStringPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}
