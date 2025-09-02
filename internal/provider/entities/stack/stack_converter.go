package stack

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	"github.com/control-monkey/controlmonkey-sdk-go/services/stack"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*stack.Stack, bool) {
	var retVal *stack.Stack

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(stack.Stack)
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

	var data stack.Data

	if deploymentBehavior, hasChanged := deploymentBehaviorConverter(plan.DeploymentBehavior, state.DeploymentBehavior, converterType); hasChanged {
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

	if runnerConfig, hasChanged := runnerConfigConverter(plan.RunnerConfig, state.RunnerConfig, converterType); hasChanged {
		data.SetRunnerConfig(runnerConfig)
		hasChanges = true
	}

	if autoSync, hasChanged := cross_models.AutoSyncConverter(plan.AutoSync, state.AutoSync, converterType); hasChanged {
		data.SetAutoSync(autoSync)
		hasChanges = true
	}

	retVal.SetData(&data)

	return retVal, hasChanges
}

func deploymentBehaviorConverter(plan *DeploymentBehaviorModel, state *DeploymentBehaviorModel, converterType commons.ConverterType) (*stack.DeploymentBehavior, bool) {
	var retVal *stack.DeploymentBehavior

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(stack.DeploymentBehavior)
	hasChanges := false

	if state == nil {
		state = new(DeploymentBehaviorModel) // dummy initialization
		hasChanges = true                    // must have changes because before is null and after is not
	}

	if plan.DeployOnPush != state.DeployOnPush {
		retVal.SetDeployOnPush(plan.DeployOnPush.ValueBoolPointer())
		hasChanges = true
	}
	if plan.WaitForApproval != state.WaitForApproval {
		retVal.SetWaitForApproval(plan.WaitForApproval.ValueBoolPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}

func vcsInfoConverter(plan *VcsInfoModel, state *VcsInfoModel, converterType commons.ConverterType) (*stack.VcsInfo, bool) {
	var retVal *stack.VcsInfo

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(stack.VcsInfo)
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

func policyConverter(plan *PolicyModel, state *PolicyModel, converterType commons.ConverterType) (*stack.Policy, bool) {
	var retVal *stack.Policy

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(stack.Policy)
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

func ttlConfigConverter(plan *TtlConfigModel, state *TtlConfigModel, converterType commons.ConverterType) (*stack.TtlConfig, bool) {
	var retVal *stack.TtlConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(stack.TtlConfig)
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

func ttlDefinitionModelConverter(plan *TtlDefinitionModel, state *TtlDefinitionModel, converterType commons.ConverterType) (*stack.TtlDefinition, bool) {
	var retVal *stack.TtlDefinition

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(stack.TtlDefinition)
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

func runnerConfigConverter(plan *RunnerConfigModel, state *RunnerConfigModel, converterType commons.ConverterType) (*stack.RunnerConfig, bool) {
	var retVal *stack.RunnerConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(stack.RunnerConfig)
	hasChanges := false

	if state == nil {
		state = new(RunnerConfigModel) // dummy initialization
		hasChanges = true              // must have changes because before is null and after is not
	}

	if plan.Mode != state.Mode {
		retVal.SetMode(plan.Mode.ValueStringPointer())
		hasChanges = true
	}

	if innerProperty, hasInnerChanges := helpers.TfListStringConverter(plan.Groups, state.Groups); hasInnerChanges {
		retVal.SetGroups(innerProperty)
		hasChanges = true
	}

	return retVal, hasChanges
}
