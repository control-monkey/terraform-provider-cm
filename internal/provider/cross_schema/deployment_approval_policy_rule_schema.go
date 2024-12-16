package cross_schema

import (
	"fmt"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var DeploymentApprovalPolicyRuleSchema = schema.ListNestedAttribute{
	MarkdownDescription: "Set up rules for approving deployment processes. At least one rule should be configured",
	Required:            true,
	Validators: []validator.List{
		listvalidator.SizeAtLeast(1),
	},
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the rule. Find supported types [here](https://docs.controlmonkey.io/controlmonkey-api/api-enumerations#deployment-approval-policy-rule-types)",
				Required:            true,
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("JSON format of the rule parameters according to the `type`. Find supported parameters [here](https://docs.controlmonkey.io/controlmonkey-api/approval-policy-rules)"),
				Optional:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
		},
	},
}
