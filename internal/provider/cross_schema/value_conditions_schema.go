package cross_schema

import (
	"fmt"

	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	cm_stringvalidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ValueConditionsSchema = schema.ListNestedAttribute{
	Optional:            true,
	MarkdownDescription: "Specify conditions for the variable value using an operator and another value. Typically used for stacks launched from templates. For more information: [ControlMonkey Docs] (https://docs.controlmonkey.io/main-concepts/variables/variable-conditions)",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"operator": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: fmt.Sprintf("Logical operators. Allowed values: %s.", helpers.EnumForDocs(cmTypes.VariableConditionOperatorTypes)),
				Validators: []validator.String{
					stringvalidator.OneOf(cmTypes.VariableConditionOperatorTypes...),
				},
			},
			"value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: fmt.Sprintf("The value associated with the operator. Input a number or string depending on the chosen operator. Use `values` field for operator of type `%s`", cmTypes.In),
				Validators:          []validator.String{cm_stringvalidators.NotBlank(), stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("values"))},
			},
			"values": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: fmt.Sprintf("A list of strings when using operator type `%s`. For other operators use `value`", cmTypes.In),
				Validators:          commons.ValidateUniqueNotEmptyListWithNoBlankValues(),
			},
		},
	},
}
