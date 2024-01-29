package namespace

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	"github.com/control-monkey/controlmonkey-sdk-go/services/namespace"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
	"reflect"
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

	if policy, hasChanged := policyConverter(plan.Policy, state.Policy, converterType); hasChanged {
		retVal.SetPolicy(policy)
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

	return retVal, hasChanges
}

func externalCredentialsConverter(plan []*ExternalCredentialsModel, state []*ExternalCredentialsModel, converterType commons.ConverterType) ([]*namespace.ExternalCredentials, bool) {
	var retVal []*namespace.ExternalCredentials
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true
		retVal = make([]*namespace.ExternalCredentials, 0)

		for _, r := range plan {
			rule := credentialsConverter(r)
			retVal = append(retVal, rule)
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

func policyConverter(plan *PolicyModel, state *PolicyModel, converterType commons.ConverterType) (*namespace.Policy, bool) {
	var retVal *namespace.Policy

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(namespace.Policy)
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

func ttlConfigConverter(plan *TtlConfigModel, state *TtlConfigModel, converterType commons.ConverterType) (*namespace.TtlConfig, bool) {
	var retVal *namespace.TtlConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(namespace.TtlConfig)
	hasChanges := false

	if state == nil {
		state = new(TtlConfigModel) // dummy initialization
		hasChanges = true           // must have changes because before is null and after is not
	}

	if defaultTtl, hasInnerChanges := ttlDefinitionModelConverter(plan.DefaultTtl, state.DefaultTtl, converterType); hasInnerChanges {
		retVal.SetDefaultTtl(defaultTtl)
		hasChanges = true
	}

	if maxTtl, hasInnerChanges := ttlDefinitionModelConverter(plan.MaxTtl, state.MaxTtl, converterType); hasInnerChanges {
		retVal.SetMaxTtl(maxTtl)
		hasChanges = true
	}
	return retVal, hasChanges
}

func ttlDefinitionModelConverter(plan *TtlDefinitionModel, state *TtlDefinitionModel, converterType commons.ConverterType) (*namespace.TtlDefinition, bool) {
	var retVal *namespace.TtlDefinition

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(namespace.TtlDefinition)
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

	return retVal, hasChanges
}
