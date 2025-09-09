package stack_discovery_configuration

import (
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID          types.String       `tfsdk:"id"`
	Name        types.String       `tfsdk:"name"`
	NamespaceId types.String       `tfsdk:"namespace_id"`
	Description types.String       `tfsdk:"description"`
	VcsPatterns []*VcsPatternModel `tfsdk:"vcs_patterns"`
	StackConfig *StackConfigModel  `tfsdk:"stack_config"`
}

type VcsPatternModel struct {
	ProviderId          types.String `tfsdk:"provider_id"`
	RepoName            types.String `tfsdk:"repo_name"`
	PathPatterns        types.List   `tfsdk:"path_patterns"`
	ExcludePathPatterns types.List   `tfsdk:"exclude_path_patterns"`
	Branch              types.String `tfsdk:"branch"`
}

type StackConfigModel struct {
	IacType                  types.String                                `tfsdk:"iac_type"`
	DeploymentBehavior       *cross_models.DeploymentBehaviorModel       `tfsdk:"deployment_behavior"`
	DeploymentApprovalPolicy *cross_models.DeploymentApprovalPolicyModel `tfsdk:"deployment_approval_policy"`
	RunTrigger               *cross_models.RunTriggerModel               `tfsdk:"run_trigger"`
	IacConfig                *cross_models.IacConfigModel                `tfsdk:"iac_config"`
	RunnerConfig             *cross_models.RunnerConfigModel             `tfsdk:"runner_config"`
	AutoSync                 *cross_models.AutoSyncModel                 `tfsdk:"auto_sync"`
}
