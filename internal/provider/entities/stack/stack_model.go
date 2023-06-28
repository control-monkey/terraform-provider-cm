package stack

import "github.com/hashicorp/terraform-plugin-framework/types"

type ResourceModel struct {
	ID                 types.String             `tfsdk:"id"`
	IacType            types.String             `tfsdk:"iac_type"`
	NamespaceId        types.String             `tfsdk:"namespace_id"`
	Name               types.String             `tfsdk:"name"`
	Description        types.String             `tfsdk:"description"`
	DeploymentBehavior *DeploymentBehaviorModel `tfsdk:"deployment_behavior"`
	VcsInfo            *VcsInfoModel            `tfsdk:"vcs_info"`
	RunTrigger         *RunTriggerModel         `tfsdk:"run_trigger"`
	IacConfig          *IacConfigModel          `tfsdk:"iac_config"`
	Policy             *PolicyModel             `tfsdk:"policy"`
}

type DeploymentBehaviorModel struct {
	DeployOnPush    types.Bool `tfsdk:"deploy_on_push"`
	WaitForApproval types.Bool `tfsdk:"wait_for_approval"`
}

type VcsInfoModel struct {
	ProviderId types.String `tfsdk:"provider_id"`
	RepoName   types.String `tfsdk:"repo_name"`
	Path       types.String `tfsdk:"path"`
	Branch     types.String `tfsdk:"branch"`
}

type RunTriggerModel struct {
	Patterns []types.String `tfsdk:"patterns"`
}

type IacConfigModel struct {
	TerraformVersion  types.String `tfsdk:"terraform_version"`
	TerragruntVersion types.String `tfsdk:"terragrunt_version"`
}

type PolicyModel struct {
	TtlConfig *TtlConfigModel `tfsdk:"ttl_config"`
}

type TtlConfigModel struct {
	Ttl *TtlDefinitionModel `tfsdk:"ttl"`
}

type TtlDefinitionModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.Int64  `tfsdk:"value"`
}

var IacTypes = []string{"terraform", "terragrunt"}
var TtlTypes = []string{"hours", "days"}
