package cross_models

import (
	sdkCrossModels "github.com/control-monkey/controlmonkey-sdk-go/services/cross_models"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
)

//region Model

type DeploymentApprovalPolicyRuleModel struct {
	Type types.String `tfsdk:"type"`
}

//endregion

//region Create/Update Converter

func DeploymentApprovalPolicyRulesConverter(plan []*DeploymentApprovalPolicyRuleModel, state []*DeploymentApprovalPolicyRuleModel, converterType commons.ConverterType) ([]*sdkCrossModels.DeploymentApprovalPolicyRule, bool) {
	var retVal []*sdkCrossModels.DeploymentApprovalPolicyRule
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true
		retVal = make([]*sdkCrossModels.DeploymentApprovalPolicyRule, 0)

		for _, r := range plan {
			if r != nil {
				rule := deploymentApprovalPolicyRuleConverter(r)
				retVal = append(retVal, rule)
			}
		}
	}

	return retVal, hasChanged
}

func deploymentApprovalPolicyRuleConverter(plan *DeploymentApprovalPolicyRuleModel) *sdkCrossModels.DeploymentApprovalPolicyRule {
	retVal := new(sdkCrossModels.DeploymentApprovalPolicyRule)

	retVal.SetType(plan.Type.ValueStringPointer())

	return retVal
}

//endregion

//region Update State After Read

func UpdateStateAfterReadDeploymentApprovalPolicyRules(deploymentApprovalPolicyRules []*sdkCrossModels.DeploymentApprovalPolicyRule) []*DeploymentApprovalPolicyRuleModel {
	var retVal []*DeploymentApprovalPolicyRuleModel

	if deploymentApprovalPolicyRules != nil {
		retVal = make([]*DeploymentApprovalPolicyRuleModel, 0)

		for _, rule := range deploymentApprovalPolicyRules {
			sr := updateStateAfterReadDeploymentApprovalPolicyRule(rule)
			retVal = append(retVal, &sr)
		}
	}

	return retVal
}

func updateStateAfterReadDeploymentApprovalPolicyRule(deploymentApprovalPolicyRule *sdkCrossModels.DeploymentApprovalPolicyRule) DeploymentApprovalPolicyRuleModel {
	var retVal DeploymentApprovalPolicyRuleModel

	retVal.Type = helpers.StringValueOrNull(deploymentApprovalPolicyRule.Type)

	return retVal
}

//endregion
