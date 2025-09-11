package stack_discovery_configuration

import (
	sdkstackdiscoveryconfig "github.com/control-monkey/controlmonkey-sdk-go/services/stack_discovery_configuration"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(res *sdkstackdiscoveryconfig.StackDiscoveryConfiguration, state *ResourceModel) {
	state.ID = helpers.StringValueOrNull(res.ID)
	state.Name = helpers.StringValueOrNull(res.Name)
	state.NamespaceId = helpers.StringValueOrNull(res.NamespaceId)
	state.Description = helpers.StringValueOrNull(res.Description)

	if res.VcsPatterns != nil {
		vcsPatterns := updateStateAfterReadVcsPatterns(res.VcsPatterns)
		state.VcsPatterns = vcsPatterns
	} else {
		state.VcsPatterns = nil
	}

	if res.StackConfig != nil {
		stackConfig := updateStateAfterReadStackConfig(res.StackConfig)
		state.StackConfig = &stackConfig
	} else {
		state.StackConfig = nil
	}
}

func updateStateAfterReadVcsPatterns(vcsPatterns []*sdkstackdiscoveryconfig.VcsPattern) []*VcsPatternModel {
	var retVal []*VcsPatternModel

	if vcsPatterns != nil {
		retVal = make([]*VcsPatternModel, 0)

		for _, p := range vcsPatterns {
			pattern := updateStateAfterReadVcsPattern(p)
			retVal = append(retVal, &pattern)
		}
	}

	return retVal
}

func updateStateAfterReadVcsPattern(pattern *sdkstackdiscoveryconfig.VcsPattern) VcsPatternModel {
	var retVal VcsPatternModel

	retVal.ProviderId = helpers.StringValueOrNull(pattern.ProviderId)
	retVal.RepoName = helpers.StringValueOrNull(pattern.RepoName)
	retVal.Branch = helpers.StringValueOrNull(pattern.Branch)

	if pattern.PathPatterns != nil {
		retVal.PathPatterns = helpers.StringPointerSliceToTfList(pattern.PathPatterns)
	} else {
		retVal.PathPatterns = types.ListNull(types.StringType)
	}

	if pattern.ExcludePathPatterns != nil {
		retVal.ExcludePathPatterns = helpers.StringPointerSliceToTfList(pattern.ExcludePathPatterns)
	} else {
		retVal.ExcludePathPatterns = types.ListNull(types.StringType)
	}

	return retVal
}

func updateStateAfterReadStackConfig(stackConfig *sdkstackdiscoveryconfig.StackConfig) StackConfigModel {
	var retVal StackConfigModel

	retVal.IacType = helpers.StringValueOrNull(stackConfig.IacType)

	if stackConfig.DeploymentBehavior != nil {
		deploymentBehavior := cross_models.UpdateStateAfterReadDeploymentBehavior(stackConfig.DeploymentBehavior)
		retVal.DeploymentBehavior = &deploymentBehavior
	}

	if stackConfig.DeploymentApprovalPolicy != nil {
		deploymentApprovalPolicy := cross_models.UpdateStateAfterReadDeploymentApprovalPolicy(stackConfig.DeploymentApprovalPolicy)
		retVal.DeploymentApprovalPolicy = &deploymentApprovalPolicy
	}

	if stackConfig.RunTrigger != nil {
		runTrigger := cross_models.UpdateStateAfterReadRunTrigger(stackConfig.RunTrigger)
		retVal.RunTrigger = &runTrigger
	}

	if stackConfig.IacConfig != nil {
		iacConfig := cross_models.UpdateStateAfterReadIacConfig(stackConfig.IacConfig)
		retVal.IacConfig = &iacConfig
	}

	if stackConfig.RunnerConfig != nil {
		runnerConfig := cross_models.UpdateStateAfterReadRunnerConfig(stackConfig.RunnerConfig)
		retVal.RunnerConfig = &runnerConfig
	}

	if stackConfig.AutoSync != nil {
		autoSync := cross_models.UpdateStateAfterReadAutoSync(stackConfig.AutoSync)
		retVal.AutoSync = &autoSync
	}

	return retVal
}
