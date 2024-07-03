package namespace

import (
	sdkNamespace "github.com/control-monkey/controlmonkey-sdk-go/services/namespace"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
)

func UpdateStateAfterRead(res *sdkNamespace.Namespace, state *ResourceModel) {
	namespace := res

	state.Name = helpers.StringValueOrNull(namespace.Name)
	state.Description = helpers.StringValueIfNotEqual(namespace.Description, "")

	if namespace.ExternalCredentials != nil {
		ec := updateStateAfterReadExternalCredentials(namespace.ExternalCredentials)
		state.ExternalCredentials = ec
	} else {
		state.ExternalCredentials = nil
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

func updateStateAfterReadIacConfig(iacConfig *sdkNamespace.IacConfig) IacConfigModel {
	var retVal IacConfigModel

	retVal.TerraformVersion = helpers.StringValueOrNull(iacConfig.TerraformVersion)
	retVal.TerragruntVersion = helpers.StringValueOrNull(iacConfig.TerragruntVersion)
	retVal.OpentofuVersion = helpers.StringValueOrNull(iacConfig.OpentofuVersion)

	return retVal
}

func updateStateAfterReadRunnerConfig(rc *sdkNamespace.RunnerConfig) RunnerConfigModel {
	var retVal RunnerConfigModel

	retVal.Mode = helpers.StringValueOrNull(rc.Mode)
	retVal.Groups = helpers.StringPointerSliceToTfList(rc.Groups)
	retVal.IsOverridable = helpers.BoolValueOrNull(rc.IsOverridable)

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

	retVal.OverrideBehavior = helpers.StringValueOrNull(deploymentApprovalPolicy.OverrideBehavior)

	return retVal
}
