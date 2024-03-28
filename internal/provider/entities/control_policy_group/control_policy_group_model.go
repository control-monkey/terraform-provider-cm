package controlPolicyGroup

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID              types.String          `tfsdk:"id"`
	Name            types.String          `tfsdk:"name"`
	Description     types.String          `tfsdk:"description"`
	ControlPolicies []*ControlPolicyModel `tfsdk:"control_policies"`
}

type ControlPolicyModel struct {
	ControlPolicyId types.String `tfsdk:"control_policy_id"`
	Severity        types.String `tfsdk:"severity"`
}
