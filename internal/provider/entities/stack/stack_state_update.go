package stack

import (
	sdkStack "github.com/control-monkey/controlmonkey-sdk-go/services/stack"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
)

func UpdateStateAfterRead(res *sdkStack.Stack, state *ResourceModel) {
	stack := res

	state.IacType = helpers.StringValueOrNull(stack.IacType)
	state.NamespaceId = helpers.StringValueOrNull(stack.NamespaceId)
	state.Name = helpers.StringValueOrNull(stack.Name)
	state.Description = helpers.StringValueIfNotEqual(stack.Description, "")

	data := stack.Data

	if data.DeploymentBehavior != nil {
		dp := updateStateAfterReadDeploymentBehavior(data.DeploymentBehavior)
		state.DeploymentBehavior = &dp
	} else {
		state.DeploymentBehavior = nil
	}

	if data.DeploymentApprovalPolicy != nil {
		dap := cross_models.UpdateStateAfterReadDeploymentApprovalPolicy(data.DeploymentApprovalPolicy)
		state.DeploymentApprovalPolicy = &dap
	} else {
		state.DeploymentApprovalPolicy = nil
	}

	if data.VcsInfo != nil {
		vcs := updateStateAfterReadVcsInfo(data.VcsInfo)
		state.VcsInfo = &vcs
	} else {
		state.VcsInfo = nil
	}

	if data.RunTrigger != nil {
		rt := cross_models.UpdateStateAfterReadRunTrigger(data.RunTrigger)
		state.RunTrigger = &rt
	} else {
		state.RunTrigger = nil
	}

	if data.IacConfig != nil {
		ic := cross_models.UpdateStateAfterReadIacConfig(data.IacConfig)
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

	if data.RunnerConfig != nil {
		rc := updateStateAfterReadRunnerConfig(data.RunnerConfig)
		state.RunnerConfig = &rc
	} else {
		state.RunnerConfig = nil
	}

	if data.AutoSync != nil {
		as := cross_models.UpdateStateAfterReadAutoSync(data.AutoSync)
		state.AutoSync = &as
	} else {
		state.AutoSync = nil
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

func updateStateAfterReadRunnerConfig(rc *sdkStack.RunnerConfig) RunnerConfigModel {
	var retVal RunnerConfigModel

	if rc != nil {
		retVal.Mode = helpers.StringValueOrNull(rc.Mode)
		retVal.Groups = helpers.StringPointerSliceToTfList(rc.Groups)
	}

	return retVal
}
