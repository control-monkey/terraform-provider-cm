package namespace

import (
	"reflect"

	"github.com/control-monkey/controlmonkey-sdk-go/services/namespace"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*namespace.Namespace, bool) {
	var retVal *namespace.Namespace

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(namespace.Namespace)
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

	if ec, hasChanged := externalCredentialsConverter(plan.ExternalCredentials, state.ExternalCredentials, converterType); hasChanged {
		retVal.SetExternalCredentials(ec)
		hasChanges = true
	}

	if iacConfig, hasChanged := iacConfigConverter(plan.IacConfig, state.IacConfig, converterType); hasChanged {
		retVal.SetIacConfig(iacConfig)
		hasChanges = true
	}

	if runnerConfig, hasChanged := runnerConfigConverter(plan.RunnerConfig, state.RunnerConfig, converterType); hasChanged {
		retVal.SetRunnerConfig(runnerConfig)
		hasChanges = true
	}

	if deploymentApprovalPolicy, hasChanged := deploymentApprovalPolicyConverter(plan.DeploymentApprovalPolicy, state.DeploymentApprovalPolicy, converterType); hasChanged {
		retVal.SetDeploymentApprovalPolicy(deploymentApprovalPolicy)
		hasChanges = true
	}

	if capabilities, hasChanged := capabilitiesConverter(plan.Capabilities, state.Capabilities, converterType); hasChanged {
		retVal.SetCapabilities(capabilities)
		hasChanges = true
	}

	return retVal, hasChanges
}

func externalCredentialsConverter(plan []*ExternalCredentialsModel, state []*ExternalCredentialsModel, converterType commons.ConverterType) ([]*namespace.ExternalCredentials, bool) {
	var retVal []*namespace.ExternalCredentials
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true

		if plan != nil {
			retVal = make([]*namespace.ExternalCredentials, 0)

			for _, r := range plan {
				rule := credentialsConverter(r)
				retVal = append(retVal, rule)
			}
		}
	}

	return retVal, hasChanged
}

func credentialsConverter(plan *ExternalCredentialsModel) *namespace.ExternalCredentials {
	retVal := new(namespace.ExternalCredentials)

	retVal.SetExternalCredentialsId(plan.ExternalCredentialsId.ValueStringPointer())
	retVal.SetType(plan.Type.ValueStringPointer())
	retVal.SetAwsProfileName(plan.AwsProfileName.ValueStringPointer())

	return retVal
}

func iacConfigConverter(plan *IacConfigModel, state *IacConfigModel, converterType commons.ConverterType) (*namespace.IacConfig, bool) {
	var retVal *namespace.IacConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(namespace.IacConfig)
	hasChanges := false

	if state == nil {
		state = new(IacConfigModel) // dummy initialization
		hasChanges = true           // must have changes because before is null and after is not
	}

	if plan.TerraformVersion != state.TerraformVersion {
		retVal.SetTerraformVersion(plan.TerraformVersion.ValueStringPointer())
		hasChanges = true
	}
	if plan.TerragruntVersion != state.TerragruntVersion {
		retVal.SetTerragruntVersion(plan.TerragruntVersion.ValueStringPointer())
		hasChanges = true
	}
	if plan.OpentofuVersion != state.OpentofuVersion {
		retVal.SetOpentofuVersion(plan.OpentofuVersion.ValueStringPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}

func runnerConfigConverter(plan *RunnerConfigModel, state *RunnerConfigModel, converterType commons.ConverterType) (*namespace.RunnerConfig, bool) {
	var retVal *namespace.RunnerConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(namespace.RunnerConfig)
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

	if plan.IsOverridable != state.IsOverridable {
		retVal.SetIsOverridable(plan.IsOverridable.ValueBoolPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}

func deploymentApprovalPolicyConverter(plan *DeploymentApprovalPolicyModel, state *DeploymentApprovalPolicyModel, converterType commons.ConverterType) (*namespace.DeploymentApprovalPolicy, bool) {
	var retVal *namespace.DeploymentApprovalPolicy

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(namespace.DeploymentApprovalPolicy)
	hasChanges := false

	if state == nil {
		state = new(DeploymentApprovalPolicyModel) // dummy initialization
		hasChanges = true                          // must have changes because before is null and after is not
	}

	if innerProperty, hasInnerChanges := cross_models.DeploymentApprovalPolicyRulesConverter(plan.Rules, state.Rules, converterType); hasInnerChanges {
		retVal.SetRules(innerProperty)
		hasChanges = true
	}

	if plan.OverrideBehavior != state.OverrideBehavior {
		retVal.SetOverrideBehavior(plan.OverrideBehavior.ValueStringPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}

func capabilitiesConverter(plan *CapabilitiesModel, state *CapabilitiesModel, converterType commons.ConverterType) (*namespace.Capabilities, bool) {
	var retVal *namespace.Capabilities

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(namespace.Capabilities)
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

func capabilityConfigConverter(plan *CapabilityConfigModel, state *CapabilityConfigModel, converterType commons.ConverterType) (*namespace.CapabilityConfig, bool) {
	var retVal *namespace.CapabilityConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(namespace.CapabilityConfig)
	hasChanges := false

	if state == nil {
		state = new(CapabilityConfigModel) // dummy initialization
		hasChanges = true                  // must have changes because before is null and after is not
	}

	if plan.Status != state.Status {
		retVal.SetStatus(plan.Status.ValueStringPointer())
		hasChanges = true
	}

	if plan.IsOverridable != state.IsOverridable {
		retVal.SetIsOverridable(plan.IsOverridable.ValueBoolPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}
