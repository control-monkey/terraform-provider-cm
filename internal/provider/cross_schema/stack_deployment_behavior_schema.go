package cross_schema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var StackDeploymentBehaviorSchema = schema.SingleNestedAttribute{
	MarkdownDescription: "The deployment behavior configuration.",
	Required:            true,
	Attributes: map[string]schema.Attribute{
		"deploy_on_push": schema.BoolAttribute{
			MarkdownDescription: "Choose whether to initiate a deployment when a push event occurs or not.",
			Required:            true,
		},
		"wait_for_approval": schema.BoolAttribute{
			MarkdownDescription: "Use `deployment_approval_policy`. Decide whether to wait for approval before proceeding with the deployment or not.",
			Optional:            true,
			DeprecationMessage:  "Attribute \"deployment_behavior.wait_for_approval\" is deprecated. Use \"deployment_approval_policy\" instead",
			Validators: []validator.Bool{
				boolvalidator.ConflictsWith(
					path.MatchRoot("deployment_approval_policy")),
			},
		},
	},
}
