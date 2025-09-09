package stack

import (
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID                       types.String                                `tfsdk:"id"`
	IacType                  types.String                                `tfsdk:"iac_type"`
	NamespaceId              types.String                                `tfsdk:"namespace_id"`
	Name                     types.String                                `tfsdk:"name"`
	Description              types.String                                `tfsdk:"description"`
	DeploymentBehavior       *cross_models.DeploymentBehaviorModel       `tfsdk:"deployment_behavior"`
	DeploymentApprovalPolicy *cross_models.DeploymentApprovalPolicyModel `tfsdk:"deployment_approval_policy"`
	VcsInfo                  *VcsInfoModel                               `tfsdk:"vcs_info"`
	RunTrigger               *cross_models.RunTriggerModel               `tfsdk:"run_trigger"`
	IacConfig                *cross_models.IacConfigModel                `tfsdk:"iac_config"`
	Policy                   *PolicyModel                                `tfsdk:"policy"`
	RunnerConfig             *cross_models.RunnerConfigModel             `tfsdk:"runner_config"`
	Capabilities             *CapabilitiesModel                          `tfsdk:"capabilities"`
	AutoSync                 *cross_models.AutoSyncModel                 `tfsdk:"auto_sync"`
}

type VcsInfoModel struct {
	ProviderId types.String `tfsdk:"provider_id"`
	RepoName   types.String `tfsdk:"repo_name"`
	Path       types.String `tfsdk:"path"`
	Branch     types.String `tfsdk:"branch"`
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

type CapabilitiesModel struct {
	DeployOnPush   *CapabilityConfigModel `tfsdk:"deploy_on_push"`
	PlanOnPr       *CapabilityConfigModel `tfsdk:"plan_on_pr"`
	DriftDetection *CapabilityConfigModel `tfsdk:"drift_detection"`
}

type CapabilityConfigModel struct {
	Status types.String `tfsdk:"status"`
}
