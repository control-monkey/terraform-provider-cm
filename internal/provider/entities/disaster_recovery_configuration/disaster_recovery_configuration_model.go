package disaster_recovery_configuration

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID             types.String         `tfsdk:"id"`
	Scope          types.String         `tfsdk:"scope"`
	CloudAccountId types.String         `tfsdk:"cloud_account_id"`
	BackupStrategy *BackupStrategyModel `tfsdk:"backup_strategy"`
}

type BackupStrategyModel struct {
	IncludeManagedResources types.Bool           `tfsdk:"include_managed_resources"`
	Mode                    types.String         `tfsdk:"mode"`
	VcsInfo                 *VcsInfoModel        `tfsdk:"vcs_info"`
	Groups                  jsontypes.Normalized `tfsdk:"groups"`
}

type VcsInfoModel struct {
	ProviderId types.String `tfsdk:"provider_id"`
	RepoName   types.String `tfsdk:"repo_name"`
	Branch     types.String `tfsdk:"branch"`
}
