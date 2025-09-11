package cross_models

import (
	"encoding/json"
	"reflect"

	sdkCrossModels "github.com/control-monkey/controlmonkey-sdk-go/services/cross_models"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//region Model

type DeploymentApprovalPolicyRuleModel struct {
	Type       types.String         `tfsdk:"type"`
	Parameters jsontypes.Normalized `tfsdk:"parameters"`
}

//endregion

//region Create/Update Converter

func DeploymentApprovalPolicyRulesConverter(plan []*DeploymentApprovalPolicyRuleModel, state []*DeploymentApprovalPolicyRuleModel, converterType commons.ConverterType) ([]*sdkCrossModels.DeploymentApprovalPolicyRule, bool) {
	var retVal []*sdkCrossModels.DeploymentApprovalPolicyRule
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true

		if plan != nil {
			retVal = make([]*sdkCrossModels.DeploymentApprovalPolicyRule, 0)

			for _, r := range plan {
				if r != nil {
					rule := deploymentApprovalPolicyRuleConverter(r)
					retVal = append(retVal, rule)
				}
			}
		}
	}

	return retVal, hasChanged
}

func deploymentApprovalPolicyRuleConverter(plan *DeploymentApprovalPolicyRuleModel) *sdkCrossModels.DeploymentApprovalPolicyRule {
	retVal := new(sdkCrossModels.DeploymentApprovalPolicyRule)

	retVal.SetType(plan.Type.ValueStringPointer())

	if plan.Parameters.IsNull() == false {
		parameters := new(map[string]any)
		plan.Parameters.Unmarshal(parameters)
		retVal.SetParameters(parameters)
	} //not sending null because not all types support parameters

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

	if deploymentApprovalPolicyRule.Parameters != nil {
		jsonSettingsString, err := json.Marshal(deploymentApprovalPolicyRule.Parameters)
		if err != nil {
			retVal.Parameters = jsontypes.NewNormalizedNull()
		} else {
			retVal.Parameters = jsontypes.NewNormalizedValue(string(jsonSettingsString))
		}
	} else {
		retVal.Parameters = jsontypes.NewNormalizedNull()
	}

	return retVal
}

//endregion
