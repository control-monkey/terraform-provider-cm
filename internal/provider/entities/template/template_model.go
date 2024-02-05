package template

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID                        types.String  `tfsdk:"id"`
	Name                      types.String  `tfsdk:"name"`
	IacType                   types.String  `tfsdk:"iac_type"`
	Description               types.String  `tfsdk:"description"`
	VcsInfo                   *VcsInfoModel `tfsdk:"vcs_info"`
	Policy                    *PolicyModel  `tfsdk:"policy"`
	SkipStateRefreshOnDestroy types.Bool    `tfsdk:"skip_state_refresh_on_destroy"`
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
	MaxTtl     *TtlDefinitionModel `tfsdk:"max_ttl"`
	DefaultTtl *TtlDefinitionModel `tfsdk:"default_ttl"`
}

type TtlDefinitionModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.Int64  `tfsdk:"value"`
}
