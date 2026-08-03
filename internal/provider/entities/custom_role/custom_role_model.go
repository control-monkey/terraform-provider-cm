package customRole

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID               types.String       `tfsdk:"id"`
	Type             types.String       `tfsdk:"type"`
	Name             types.String       `tfsdk:"name"`
	Description      types.String       `tfsdk:"description"`
	Permissions      []*PermissionModel `tfsdk:"permissions"`
	StackRestriction types.String       `tfsdk:"stack_restriction"`
}

type PermissionModel struct {
	Name         types.String        `tfsdk:"name"`
	Names        []types.String      `tfsdk:"names"`
	Restrictions []*RestrictionModel `tfsdk:"restrictions"`
}

type RestrictionModel struct {
	CloudProvider  types.String `tfsdk:"cloud_provider"`
	CloudAccountId types.String `tfsdk:"cloud_account_id"`
	CmResourceName types.String `tfsdk:"cm_resource_name"`
}
