package variable

import (
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID              types.String                   `tfsdk:"id"`
	Scope           types.String                   `tfsdk:"scope"`
	ScopeId         types.String                   `tfsdk:"scope_id"`
	Key             types.String                   `tfsdk:"key"`
	Type            types.String                   `tfsdk:"type"`
	Value           types.String                   `tfsdk:"value"`
	DisplayName     types.String                   `tfsdk:"display_name"`
	IsSensitive     types.Bool                     `tfsdk:"is_sensitive"`
	IsOverridable   types.Bool                     `tfsdk:"is_overridable"`
	IsRequired      types.Bool                     `tfsdk:"is_required"`
	Description     types.String                   `tfsdk:"description"`
	ValueConditions []*cross_models.ConditionModel `tfsdk:"value_conditions"`
}
