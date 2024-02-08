package controlPolicyGroupMappings

import (
	"fmt"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type ResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	ControlPolicyGroupId types.String   `tfsdk:"control_policy_group_id"`
	Targets              []*TargetModel `tfsdk:"targets"`
}

type TargetModel struct {
	TargetId             types.String                `tfsdk:"target_id"`
	TargetType           types.String                `tfsdk:"target_type"`
	EnforcementLevel     types.String                `tfsdk:"enforcement_level"`
	OverrideEnforcements []*OverrideEnforcementModel `tfsdk:"override_enforcements"`
}

type OverrideEnforcementModel struct {
	ControlPolicyId  types.String `tfsdk:"control_policy_id"`
	EnforcementLevel types.String `tfsdk:"enforcement_level"`
	StackIds         types.List   `tfsdk:"stack_ids"`
}

func (e *TargetModel) Hash() string {
	retVal := ""

	if e.TargetId.IsNull() == false {
		retVal += fmt.Sprintf("TargetId:%s:", e.TargetId.ValueString())
	}
	if e.TargetType.IsNull() == false {
		retVal += fmt.Sprintf("TargetType:%s:", e.TargetType.ValueString())
	}
	if e.EnforcementLevel.IsNull() == false {
		retVal += fmt.Sprintf("EnforcementLevel:%s:", e.EnforcementLevel.ValueString())
	}

	for _, o := range e.OverrideEnforcements {
		retVal += o.Hash()
	}

	return retVal
}

func (e *TargetModel) GetBlockIdentifier() string {
	retVal := ""

	if helpers.IsKnown(e.TargetId) && helpers.IsKnown(e.TargetType) {
		retVal = fmt.Sprintf("TargetId:%s:TargetType:%s", e.TargetId.ValueString(), e.TargetType.ValueString())

	}

	return retVal
}

func (e *OverrideEnforcementModel) Hash() string {
	retVal := ""

	if e.ControlPolicyId.IsNull() == false {
		retVal += fmt.Sprintf("ControlPolicyId:%s:", e.ControlPolicyId.ValueString())
	}
	if e.EnforcementLevel.IsNull() == false {
		retVal += fmt.Sprintf("EnforcementLevel:%s:", e.EnforcementLevel.ValueString())
	}
	if e.StackIds.IsNull() == false {
		retVal += fmt.Sprintf("StackIds:%s:", e.StackIds.String())
	}

	return retVal
}

func CleanTargetIdentifier(s string) string {
	split := strings.Split(s, ":")
	return fmt.Sprintf("%s %s", split[3], split[1])
}
