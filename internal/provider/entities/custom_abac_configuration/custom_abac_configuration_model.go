package customAbacConfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID           types.String `tfsdk:"id"`
	CustomAbacId types.String `tfsdk:"custom_abac_id"`
	Name         types.String `tfsdk:"name"`
	Roles        []*RoleModel `tfsdk:"roles"`
}

type RoleModel struct {
	OrgId   types.String `tfsdk:"org_id"`
	OrgRole types.String `tfsdk:"org_role"`
	TeamIds types.List   `tfsdk:"team_ids"`
}
