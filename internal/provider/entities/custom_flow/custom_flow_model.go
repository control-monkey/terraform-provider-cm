package custom_flow

import "github.com/hashicorp/terraform-plugin-framework/types"

type ResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	FlowYaml  types.String `tfsdk:"flow_yaml"`
	IsEnabled types.Bool   `tfsdk:"is_enabled"`
}
