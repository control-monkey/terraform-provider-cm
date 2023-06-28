package stack

import (
	sdkStack "github.com/control-monkey/controlmonkey-sdk-go/services/stack"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res *sdkStack.ReadStackOutput, state *ResourceModel) {
	stack := res.Stack

	state.IacType = helpers.StringValueOrNull(stack.IacType)
	state.NamespaceId = helpers.StringValueOrNull(stack.NamespaceId)
	state.Name = helpers.StringValueOrNull(stack.Name)
	state.Description = helpers.StringValueOrNull(stack.Description)

	data := stack.Data

	if data.DeploymentBehavior != nil {
		dp := updateStateAfterReadDeploymentBehavior(data.DeploymentBehavior)
		state.DeploymentBehavior = &dp
	} else {
		state.DeploymentBehavior = nil
	}

	if data.VcsInfo != nil {
		vcs := updateStateAfterReadVcsInfo(data.VcsInfo)
		state.VcsInfo = &vcs
	} else {
		state.VcsInfo = nil
	}

	if data.RunTrigger != nil {
		rt := updateStateAfterReadRunTrigger(data.RunTrigger)
		state.RunTrigger = &rt
	} else {
		state.RunTrigger = nil
	}

	if data.IacConfig != nil {
		ic := updateStateAfterReadIacConfig(data.IacConfig)
		state.IacConfig = &ic
	} else {
		state.IacConfig = nil
	}

	if data.Policy != nil {
		policy := updateStateAfterReadPolicy(data.Policy)
		state.Policy = &policy
	} else {
		state.Policy = nil
	}
}

func updateStateAfterReadDeploymentBehavior(deploymentBehavior *sdkStack.DeploymentBehavior) DeploymentBehaviorModel {
	var retVal DeploymentBehaviorModel

	retVal.DeployOnPush = helpers.BoolValueOrNull(deploymentBehavior.DeployOnPush)
	retVal.WaitForApproval = helpers.BoolValueOrNull(deploymentBehavior.WaitForApproval)

	return retVal
}

func updateStateAfterReadVcsInfo(vcsInfo *sdkStack.VcsInfo) VcsInfoModel {
	var retVal VcsInfoModel

	retVal.ProviderId = helpers.StringValueOrNull(vcsInfo.ProviderId)
	retVal.RepoName = helpers.StringValueOrNull(vcsInfo.RepoName)
	retVal.Path = helpers.StringValueOrNull(vcsInfo.Path)
	retVal.Branch = helpers.StringValueOrNull(vcsInfo.Branch)

	return retVal
}

func updateStateAfterReadRunTrigger(runTrigger *sdkStack.RunTrigger) RunTriggerModel {
	var retVal RunTriggerModel

	retVal.Patterns = helpers.StringSliceOrNull(runTrigger.Patterns)
	return retVal
}

func updateStateAfterReadIacConfig(iacConfig *sdkStack.IacConfig) IacConfigModel {
	var retVal IacConfigModel

	retVal.TerraformVersion = helpers.StringValueOrNull(iacConfig.TerraformVersion)
	retVal.TerragruntVersion = helpers.StringValueOrNull(iacConfig.TerragruntVersion)

	return retVal
}

func updateStateAfterReadPolicy(policy *sdkStack.Policy) PolicyModel {
	var retVal PolicyModel

	if policy.TtlConfig != nil {
		ttlConfig := updateStateAfterReadTtlConfig(policy.TtlConfig)
		retVal.TtlConfig = &ttlConfig
	} else {
		retVal.TtlConfig = nil
	}

	return retVal
}

func updateStateAfterReadTtlConfig(ttlConfig *sdkStack.TtlConfig) TtlConfigModel {
	var retVal TtlConfigModel

	if ttlConfig.Ttl != nil {
		ttl := updateStateAfterReadTtlDefinition(ttlConfig.Ttl)
		retVal.Ttl = &ttl
	} else {
		retVal.Ttl = nil
	}

	return retVal
}

func updateStateAfterReadTtlDefinition(ttl *sdkStack.TtlDefinition) TtlDefinitionModel {
	var retVal TtlDefinitionModel

	retVal.Type = helpers.StringValueOrNull(ttl.Type)
	retVal.Value = helpers.Int64ValueOrNull(ttl.Value)

	return retVal
}
