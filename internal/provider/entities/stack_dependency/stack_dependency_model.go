package stack_dependency

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID               types.String      `tfsdk:"id"`
	StackId          types.String      `tfsdk:"stack_id"`
	DependsOnStackId types.String      `tfsdk:"depends_on_stack_id"`
	TriggerOption    types.String      `tfsdk:"trigger_option"`
	References       []*ReferenceModel `tfsdk:"references"`
}

type ReferenceModel struct {
	OutputOfStackToDependOn types.String `tfsdk:"output_of_stack_to_depend_on"`
	InputForStack           types.String `tfsdk:"input_for_stack"`
	IncludeSensitiveOutput  types.Bool   `tfsdk:"include_sensitive_output"`
}
