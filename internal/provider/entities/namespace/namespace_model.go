package namespace

import (
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID                       types.String                   `tfsdk:"id"`
	Name                     types.String                   `tfsdk:"name"`
	Description              types.String                   `tfsdk:"description"`
	ExternalCredentials      []*ExternalCredentialsModel    `tfsdk:"external_credentials"`
	Policy                   *PolicyModel                   `tfsdk:"policy"`
	IacConfig                *IacConfigModel                `tfsdk:"iac_config"`
	RunnerConfig             *RunnerConfigModel             `tfsdk:"runner_config"`
	DeploymentApprovalPolicy *DeploymentApprovalPolicyModel `tfsdk:"deployment_approval_policy"`
}

type ExternalCredentialsModel struct {
	Type                  types.String `tfsdk:"type"`
	ExternalCredentialsId types.String `tfsdk:"external_credentials_id"`
	AwsProfileName        types.String `tfsdk:"aws_profile_name"`
}

type PolicyModel struct {
	TtlConfig *TtlConfigModel `tfsdk:"ttl_config"`
}

type TtlConfigModel struct {
	MaxTtl     *TtlDefinitionModel `tfsdk:"max_ttl"`
	DefaultTtl *TtlDefinitionModel `tfsdk:"default_ttl"`
}

type TtlDefinitionModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.Int64  `tfsdk:"value"`
}

type IacConfigModel struct {
	TerraformVersion  types.String `tfsdk:"terraform_version"`
	TerragruntVersion types.String `tfsdk:"terragrunt_version"`
	OpentofuVersion   types.String `tfsdk:"opentofu_version"`
}

type RunnerConfigModel struct {
	Mode          types.String   `tfsdk:"mode"`
	Groups        []types.String `tfsdk:"groups"`
	IsOverridable types.Bool     `tfsdk:"is_overridable"`
}

type DeploymentApprovalPolicyModel struct {
	Rules            []*cross_models.DeploymentApprovalPolicyRuleModel `tfsdk:"rules"`
	OverrideBehavior types.String                                      `tfsdk:"override_behavior"`
}
