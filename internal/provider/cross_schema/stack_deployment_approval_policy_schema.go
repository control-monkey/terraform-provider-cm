package cross_schema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var StackDeploymentApprovalPolicySchema = schema.SingleNestedAttribute{
	MarkdownDescription: "Set up requirements to approve a deployment",
	Optional:            true,
	Attributes: map[string]schema.Attribute{
		"rules": DeploymentApprovalPolicyRuleSchema,
	},
}
