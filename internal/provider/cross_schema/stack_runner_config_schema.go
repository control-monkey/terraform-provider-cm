package cross_schema

import (
	"fmt"

	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var StackRunnerConfigSchema = schema.SingleNestedAttribute{
	MarkdownDescription: "Configure the runner settings to specify whether ControlMonkey manages the runner or it is self-hosted.",
	Optional:            true,
	Attributes: map[string]schema.Attribute{
		"mode": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The runner mode. Allowed values: %s.", helpers.EnumForDocs(cmTypes.RunnerConfigModeTypes)),
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf(cmTypes.RunnerConfigModeTypes...),
			},
		},
		"groups": schema.ListAttribute{
			MarkdownDescription: fmt.Sprintf("In case that `mode` is `%s`, groups must contain at least one runners group. If `mode` is `%s`, this field must not be configured.", cmTypes.SelfHosted, cmTypes.Managed),
			ElementType:         types.StringType,
			Optional:            true,
			// Validation in ValidateConfig
		},
	},
}
