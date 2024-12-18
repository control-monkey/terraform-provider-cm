package blueprint

import (
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID                               types.String                `tfsdk:"id"`
	Name                             types.String                `tfsdk:"name"`
	Description                      types.String                `tfsdk:"description"`
	BlueprintVcsInfo                 *VcsInfoModel               `tfsdk:"blueprint_vcs_info"`
	StackConfiguration               *StackConfigurationModel    `tfsdk:"stack_configuration"`
	SubstituteParameters             []*SubstituteParameterModel `tfsdk:"substitute_parameters"`
	SkipPlanOnStackInitialization    types.Bool                  `tfsdk:"skip_plan_on_stack_initialization"`
	AutoApproveApplyOnInitialization types.Bool                  `tfsdk:"auto_approve_apply_on_initialization"`
}

type VcsInfoModel struct {
	ProviderId types.String `tfsdk:"provider_id"`
	RepoName   types.String `tfsdk:"repo_name"`
	Path       types.String `tfsdk:"path"`
	Branch     types.String `tfsdk:"branch"`
}

type StackConfigurationModel struct {
	NamePattern              types.String                                `tfsdk:"name_pattern"`
	IacType                  types.String                                `tfsdk:"iac_type"`
	VcsInfoWithPatterns      *StackVcsInfoWithPatternsModel              `tfsdk:"vcs_info_with_patterns"`
	DeploymentApprovalPolicy *cross_models.DeploymentApprovalPolicyModel `tfsdk:"deployment_approval_policy"`
}

type StackVcsInfoWithPatternsModel struct {
	ProviderId    types.String `tfsdk:"provider_id"`
	RepoName      types.String `tfsdk:"repo_name"`
	PathPattern   types.String `tfsdk:"path_pattern"`
	BranchPattern types.String `tfsdk:"branch_pattern"`
}

type SubstituteParameterModel struct {
	Key             types.String                   `tfsdk:"key"`
	Description     types.String                   `tfsdk:"description"`
	ValueConditions []*cross_models.ConditionModel `tfsdk:"value_conditions"`
}
