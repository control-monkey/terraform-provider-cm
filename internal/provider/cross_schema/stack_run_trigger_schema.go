package cross_schema

import (
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var RunTriggerSchema = schema.SingleNestedAttribute{
	MarkdownDescription: "Glob patterns to specify additional paths that should trigger a stack run.",
	Optional:            true,
	Attributes: map[string]schema.Attribute{
		"patterns": schema.ListAttribute{
			MarkdownDescription: "Patterns that trigger a stack run.",
			ElementType:         types.StringType,
			Optional:            true,
			Validators:          commons.ValidateUniqueNotEmptyListWithNoBlankValues(),
		},
		"exclude_patterns": schema.ListAttribute{
			MarkdownDescription: "Patterns that will not trigger a stack run.",
			ElementType:         types.StringType,
			Optional:            true,
			Validators:          commons.ValidateUniqueNotEmptyListWithNoBlankValues(),
		},
	},
}
