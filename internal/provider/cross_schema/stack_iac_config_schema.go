package cross_schema

import (
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var IacConfigSchema = schema.SingleNestedAttribute{
	MarkdownDescription: "IaC configuration.",
	Optional:            true,
	Attributes: map[string]schema.Attribute{
		"terraform_version": schema.StringAttribute{
			MarkdownDescription: "the Terraform version that will be used for terraform operations.",
			Optional:            true,
		},
		"terragrunt_version": schema.StringAttribute{
			MarkdownDescription: "the Terragrunt version that will be used for terragrunt operations.",
			Optional:            true,
		},
		"opentofu_version": schema.StringAttribute{
			MarkdownDescription: "the OpenTofu version that will be used for tofu operations.",
			Optional:            true,
		},
		"is_terragrunt_run_all": schema.BoolAttribute{
			MarkdownDescription: "When using terragrunt, as long as this field is set to `True`, this field will execute \"run-all\" commands on multiple modules for init/plan/apply",
			Optional:            true,
		},
		"var_files": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "Custom variable files to pass on to Terraform. For more information: [ControlMonkey Docs](https://docs.controlmonkey.io/main-concepts/stack/stack-settings#var-files)",
			Validators:          commons.ValidateUniqueNotEmptyListWithNoBlankValues(),
		},
	},
}
