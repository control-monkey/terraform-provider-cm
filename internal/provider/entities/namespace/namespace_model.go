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
	IacConfig                *IacConfigModel                `tfsdk:"iac_config"`
	RunnerConfig             *RunnerConfigModel             `tfsdk:"runner_config"`
	DeploymentApprovalPolicy *DeploymentApprovalPolicyModel `tfsdk:"deployment_approval_policy"`
	Capabilities             *CapabilitiesModel             `tfsdk:"capabilities"`
}

type ExternalCredentialsModel struct {
	Type                  types.String `tfsdk:"type"`
	ExternalCredentialsId types.String `tfsdk:"external_credentials_id"`
	AwsProfileName        types.String `tfsdk:"aws_profile_name"`
}

type IacConfigModel struct {
	TerraformVersion  types.String `tfsdk:"terraform_version"`
	TerragruntVersion types.String `tfsdk:"terragrunt_version"`
	OpentofuVersion   types.String `tfsdk:"opentofu_version"`
}

type RunnerConfigModel struct {
	Mode          types.String `tfsdk:"mode"`
	Groups        types.List   `tfsdk:"groups"`
	IsOverridable types.Bool   `tfsdk:"is_overridable"`
}

type DeploymentApprovalPolicyModel struct {
	Rules            []*cross_models.DeploymentApprovalPolicyRuleModel `tfsdk:"rules"`
	OverrideBehavior types.String                                      `tfsdk:"override_behavior"`
}

type CapabilitiesModel struct {
	DeployOnPush   *CapabilityConfigModel `tfsdk:"deploy_on_push"`
	PlanOnPr       *CapabilityConfigModel `tfsdk:"plan_on_pr"`
	DriftDetection *CapabilityConfigModel `tfsdk:"drift_detection"`
}

type CapabilityConfigModel struct {
	Status        types.String `tfsdk:"status"`
	IsOverridable types.Bool   `tfsdk:"is_overridable"`
}
