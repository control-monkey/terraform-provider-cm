package organization

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	ImportID = "org-config"
)

type ResourceModel struct {
	ID                    types.String                 `tfsdk:"id"`
	IacConfig             *IacConfigModel              `tfsdk:"iac_config"`
	S3StateFilesLocations []*S3StateFilesLocationModel `tfsdk:"s3_state_files_locations"`
	RunnerConfig          *RunnerConfigModel           `tfsdk:"runner_config"`
	SuppressedResources   *SuppressedResourcesModel    `tfsdk:"suppressed_resources"`
	ReportConfigurations  []*ReportConfigurationModel  `tfsdk:"report_configurations"`
}

type IacConfigModel struct {
	TerraformVersion  types.String `tfsdk:"terraform_version"`
	TerragruntVersion types.String `tfsdk:"terragrunt_version"`
	OpentofuVersion   types.String `tfsdk:"opentofu_version"`
}

type S3StateFilesLocationModel struct {
	BucketName   types.String `tfsdk:"bucket_name"`
	BucketRegion types.String `tfsdk:"bucket_region"`
	AwsAccountId types.String `tfsdk:"aws_account_id"`
}

type RunnerConfigModel struct {
	Mode          types.String `tfsdk:"mode"`
	Groups        types.List   `tfsdk:"groups"`
	IsOverridable types.Bool   `tfsdk:"is_overridable"`
}

type SuppressedResourcesModel struct {
	ManagedByTags []*TagPropertiesModel `tfsdk:"managed_by_tags"`
}

type TagPropertiesModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type ReportConfigurationModel struct {
	Type       types.String           `tfsdk:"type"`
	Recipients *ReportRecipientsModel `tfsdk:"recipients"`
	Enabled    types.Bool             `tfsdk:"enabled"`
}

type ReportRecipientsModel struct {
	AllAdmins               types.Bool `tfsdk:"all_admins"`
	EmailAddresses          types.List `tfsdk:"email_addresses"`
	EmailAddressesToExclude types.List `tfsdk:"email_addresses_to_exclude"`
}
