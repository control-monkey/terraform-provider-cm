package namespace

import (
	sdkNamespace "github.com/control-monkey/controlmonkey-sdk-go/services/namespace"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
)

func UpdateStateAfterRead(res *sdkNamespace.ReadNamespaceOutput, state *ResourceModel) {
	namespace := res.Namespace

	state.Name = helpers.StringValueOrNull(namespace.Name)
	state.Description = helpers.StringValueOrNull(namespace.Description)

	if namespace.ExternalCredentials != nil {
		ec := updateStateAfterReadExternalCredentials(namespace.ExternalCredentials)
		state.ExternalCredentials = ec
	} else {
		state.ExternalCredentials = nil
	}

	if namespace.Policy != nil {
		policy := updateStateAfterReadPolicy(namespace.Policy)
		state.Policy = &policy
	} else {
		state.Policy = nil
	}

	if namespace.IacConfig != nil {
		ic := updateStateAfterReadIacConfig(namespace.IacConfig)
		state.IacConfig = &ic
	} else {
		state.IacConfig = nil
	}

	if namespace.RunnerConfig != nil {
		rc := updateStateAfterReadRunnerConfig(namespace.RunnerConfig)
		state.RunnerConfig = &rc
	} else {
		state.RunnerConfig = nil
	}

	if namespace.DeploymentApprovalPolicy != nil {
		dap := updateStateAfterReadDeploymentApprovalPolicy(namespace.DeploymentApprovalPolicy)
		state.DeploymentApprovalPolicy = &dap
	} else {
		state.DeploymentApprovalPolicy = nil
	}

}

func updateStateAfterReadExternalCredentials(externalCredentials []*sdkNamespace.ExternalCredentials) []*ExternalCredentialsModel {
	var retVal []*ExternalCredentialsModel

	if externalCredentials != nil {
		retVal = make([]*ExternalCredentialsModel, 0)

		for _, rule := range externalCredentials {
			sr := updateStateAfterReadCredentials(rule)
			retVal = append(retVal, &sr)
		}
	}

	return retVal
}

func updateStateAfterReadCredentials(credentials *sdkNamespace.ExternalCredentials) ExternalCredentialsModel {
	var retVal ExternalCredentialsModel

	retVal.ExternalCredentialsId = helpers.StringValueOrNull(credentials.ExternalCredentialsId)
	retVal.Type = helpers.StringValueOrNull(credentials.Type)
	retVal.AwsProfileName = helpers.StringValueOrNull(credentials.AwsProfileName)

	return retVal
}

func updateStateAfterReadPolicy(policy *sdkNamespace.Policy) PolicyModel {
	var retVal PolicyModel

	if policy.TtlConfig != nil {
		ttlConfig := updateStateAfterReadTtlConfig(policy.TtlConfig)
		retVal.TtlConfig = &ttlConfig
	} else {
		retVal.TtlConfig = nil
	}

	return retVal
}

func updateStateAfterReadTtlConfig(ttlConfig *sdkNamespace.TtlConfig) TtlConfigModel {
	var retVal TtlConfigModel

	if ttlConfig.DefaultTtl != nil {
		dTtl := updateStateAfterReadTtlDefinition(ttlConfig.DefaultTtl)
		retVal.DefaultTtl = &dTtl
	} else {
		retVal.DefaultTtl = nil
	}

	if ttlConfig.MaxTtl != nil {
		mTtl := updateStateAfterReadTtlDefinition(ttlConfig.MaxTtl)
		retVal.MaxTtl = &mTtl
	} else {
		retVal.MaxTtl = nil
	}

	return retVal
}

func updateStateAfterReadTtlDefinition(ttl *sdkNamespace.TtlDefinition) TtlDefinitionModel {
	var retVal TtlDefinitionModel

	retVal.Type = helpers.StringValueOrNull(ttl.Type)
	retVal.Value = helpers.Int64ValueOrNull(ttl.Value)

	return retVal
}

func updateStateAfterReadIacConfig(iacConfig *sdkNamespace.IacConfig) IacConfigModel {
	var retVal IacConfigModel

	retVal.TerraformVersion = helpers.StringValueOrNull(iacConfig.TerraformVersion)
	retVal.TerragruntVersion = helpers.StringValueOrNull(iacConfig.TerragruntVersion)

	return retVal
}

func updateStateAfterReadRunnerConfig(rc *sdkNamespace.RunnerConfig) RunnerConfigModel {
	var retVal RunnerConfigModel

	if rc != nil {
		retVal.Mode = helpers.StringValueOrNull(rc.Mode)
		retVal.Groups = helpers.StringSliceOrNull(rc.Groups)
	}

	return retVal
}

func updateStateAfterReadDeploymentApprovalPolicy(deploymentApprovalPolicy *sdkNamespace.DeploymentApprovalPolicy) DeploymentApprovalPolicyModel {
	var retVal DeploymentApprovalPolicyModel

	if deploymentApprovalPolicy.Rules != nil {
		rs := cross_models.UpdateStateAfterReadDeploymentApprovalPolicyRules(deploymentApprovalPolicy.Rules)
		retVal.Rules = rs
	} else {
		retVal.Rules = nil
	}

	return retVal
}
