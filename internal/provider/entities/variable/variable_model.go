package variable

import "github.com/hashicorp/terraform-plugin-framework/types"

type ResourceModel struct {
	ID              types.String      `tfsdk:"id"`
	Scope           types.String      `tfsdk:"scope"`
	ScopeId         types.String      `tfsdk:"scope_id"`
	Key             types.String      `tfsdk:"key"`
	Type            types.String      `tfsdk:"type"`
	Value           types.String      `tfsdk:"value"`
	DisplayName     types.String      `tfsdk:"display_name"`
	IsSensitive     types.Bool        `tfsdk:"is_sensitive"`
	IsOverridable   types.Bool        `tfsdk:"is_overridable"`
	IsRequired      types.Bool        `tfsdk:"is_required"`
	Description     types.String      `tfsdk:"description"`
	ValueConditions []*ConditionModel `tfsdk:"value_conditions"`
}

type ConditionModel struct {
	Operator types.String `tfsdk:"operator"`
	Value    types.String `tfsdk:"value"`
	Values   types.List   `tfsdk:"values"`
}
