package customRole

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID               types.String       `tfsdk:"id"`
	Name             types.String       `tfsdk:"name"`
	Description      types.String       `tfsdk:"description"`
	Permissions      []*PermissionModel `tfsdk:"permissions"`
	StackRestriction types.String       `tfsdk:"stack_restriction"`
}

type PermissionModel struct {
	Name types.String `tfsdk:"name"`
}
