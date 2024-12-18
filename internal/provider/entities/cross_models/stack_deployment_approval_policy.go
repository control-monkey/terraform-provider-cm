package cross_models

import (
	sdkCrossModels "github.com/control-monkey/controlmonkey-sdk-go/services/cross_models"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
)

//region Model

type DeploymentApprovalPolicyModel struct {
	Rules []*DeploymentApprovalPolicyRuleModel `tfsdk:"rules"`
}

//endregion

//region Create/Update Converter

func DeploymentApprovalPolicyConverter(plan *DeploymentApprovalPolicyModel, state *DeploymentApprovalPolicyModel, converterType commons.ConverterType) (*sdkCrossModels.DeploymentApprovalPolicy, bool) {
	var retVal *sdkCrossModels.DeploymentApprovalPolicy

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(sdkCrossModels.DeploymentApprovalPolicy)
	hasChanges := false

	if state == nil {
		state = new(DeploymentApprovalPolicyModel) // dummy initialization
		hasChanges = true                          // must have changes because before is null and after is not
	}

	if innerProperty, hasInnerChanges := DeploymentApprovalPolicyRulesConverter(plan.Rules, state.Rules, converterType); hasInnerChanges {
		retVal.SetRules(innerProperty)
		hasChanges = true
	}

	return retVal, hasChanges
}

//endregion

//region Update State After Read

func UpdateStateAfterReadDeploymentApprovalPolicy(deploymentApprovalPolicy *sdkCrossModels.DeploymentApprovalPolicy) DeploymentApprovalPolicyModel {
	var retVal DeploymentApprovalPolicyModel

	if deploymentApprovalPolicy.Rules != nil {
		rs := UpdateStateAfterReadDeploymentApprovalPolicyRules(deploymentApprovalPolicy.Rules)
		retVal.Rules = rs
	} else {
		retVal.Rules = nil
	}

	return retVal
}

//endregion
