package template

import (
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID                        types.String                    `tfsdk:"id"`
	Name                      types.String                    `tfsdk:"name"`
	IacType                   types.String                    `tfsdk:"iac_type"`
	Description               types.String                    `tfsdk:"description"`
	VcsInfo                   *VcsInfoModel                   `tfsdk:"vcs_info"`
	Policy                    *PolicyModel                    `tfsdk:"policy"`
	SkipStateRefreshOnDestroy types.Bool                      `tfsdk:"skip_state_refresh_on_destroy"`
	IacConfig                 *IacConfigModel                 `tfsdk:"iac_config"`
	RunnerConfig              *cross_models.RunnerConfigModel `tfsdk:"runner_config"`
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

type IacConfigModel struct {
	TerraformVersion  types.String `tfsdk:"terraform_version"`
	TerragruntVersion types.String `tfsdk:"terragrunt_version"`
	OpentofuVersion   types.String `tfsdk:"opentofu_version"`
}
