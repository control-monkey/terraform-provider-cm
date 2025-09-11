package stack_discovery_configuration

import (
	"reflect"

	sdkstackdiscoveryconfig "github.com/control-monkey/controlmonkey-sdk-go/services/stack_discovery_configuration"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*sdkstackdiscoveryconfig.StackDiscoveryConfiguration, bool) {
	var retVal *sdkstackdiscoveryconfig.StackDiscoveryConfiguration

	if plan == nil {
		if state == nil {
			return nil, false
		} else {
			return nil, true
		}
	}

	retVal = new(sdkstackdiscoveryconfig.StackDiscoveryConfiguration)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel)
		hasChanges = true
	}

	if plan.Name != state.Name {
		retVal.SetName(plan.Name.ValueStringPointer())
		hasChanges = true
	}
	if plan.NamespaceId != state.NamespaceId {
		retVal.SetNamespaceId(plan.NamespaceId.ValueStringPointer())
		hasChanges = true
	}
	if plan.Description != state.Description {
		retVal.SetDescription(plan.Description.ValueStringPointer())
		hasChanges = true
	}

	if vcsPatterns, changed := vcsPatternsConverter(plan.VcsPatterns, state.VcsPatterns, converterType); changed {
		retVal.SetVcsPatterns(vcsPatterns)
		hasChanges = true
	}

	if stackConfig, changed := stackConfigConverter(plan.StackConfig, state.StackConfig, converterType); changed {
		retVal.SetStackConfig(stackConfig)
		hasChanges = true
	}

	return retVal, hasChanges
}

func vcsPatternsConverter(plan []*VcsPatternModel, state []*VcsPatternModel, converterType commons.ConverterType) ([]*sdkstackdiscoveryconfig.VcsPattern, bool) {
	var retVal []*sdkstackdiscoveryconfig.VcsPattern
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true

		if plan != nil {
			retVal = make([]*sdkstackdiscoveryconfig.VcsPattern, 0)

			for _, p := range plan {
				pattern := vcsPatternConverter(p)
				retVal = append(retVal, pattern)
			}
		}
	}

	return retVal, hasChanged
}

func vcsPatternConverter(plan *VcsPatternModel) *sdkstackdiscoveryconfig.VcsPattern {
	retVal := new(sdkstackdiscoveryconfig.VcsPattern)

	retVal.SetProviderId(plan.ProviderId.ValueStringPointer())
	retVal.SetRepoName(plan.RepoName.ValueStringPointer())
	retVal.SetBranch(plan.Branch.ValueStringPointer())

	if helpers.IsKnown(plan.PathPatterns) {
		pathPatterns := helpers.TfListToStringSlice(plan.PathPatterns)
		retVal.SetPathPatterns(pathPatterns)
	}

	if helpers.IsKnown(plan.ExcludePathPatterns) {
		excludePathPatterns := helpers.TfListToStringSlice(plan.ExcludePathPatterns)
		retVal.SetExcludePathPatterns(excludePathPatterns)
	}

	return retVal
}

func stackConfigConverter(plan *StackConfigModel, state *StackConfigModel, converterType commons.ConverterType) (*sdkstackdiscoveryconfig.StackConfig, bool) {
	var retVal *sdkstackdiscoveryconfig.StackConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(sdkstackdiscoveryconfig.StackConfig)
	hasChanges := false

	if state == nil {
		state = new(StackConfigModel) // dummy initialization
		hasChanges = true             // must have changes because before is null and after is not
	}

	if plan.IacType != state.IacType {
		retVal.SetIacType(plan.IacType.ValueStringPointer())
		hasChanges = true
	}

	if deploymentBehavior, hasChanged := cross_models.DeploymentBehaviorConverter(plan.DeploymentBehavior, state.DeploymentBehavior, converterType); hasChanged {
		retVal.SetDeploymentBehavior(deploymentBehavior)
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

	if runnerConfig, hasChanged := cross_models.RunnerConfigConverter(plan.RunnerConfig, state.RunnerConfig, converterType); hasChanged {
		retVal.SetRunnerConfig(runnerConfig)
		hasChanges = true
	}

	if autoSync, hasChanged := cross_models.AutoSyncConverter(plan.AutoSync, state.AutoSync, converterType); hasChanged {
		retVal.SetAutoSync(autoSync)
		hasChanges = true
	}

	return retVal, hasChanges
}
