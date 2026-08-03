package customFlowMappings

import (
	"fmt"
	"strings"

	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID           types.String   `tfsdk:"id"`
	CustomFlowId types.String   `tfsdk:"custom_flow_id"`
	Targets      []*TargetModel `tfsdk:"targets"`
}

type TargetModel struct {
	TargetId   types.String `tfsdk:"target_id"`
	TargetType types.String `tfsdk:"target_type"`
}

// Hash covers every field: there are no mutable fields, so every change is a delete + create.
func (e *TargetModel) Hash() string {
	retVal := ""

	if e.TargetId.IsNull() == false {
		retVal += fmt.Sprintf("TargetId:%s:", e.TargetId.ValueString())
	}
	if e.TargetType.IsNull() == false {
		retVal += fmt.Sprintf("TargetType:%s:", e.TargetType.ValueString())
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

func CleanIdentifier(s string) string {
	split := strings.Split(s, ":")
	return fmt.Sprintf("%s %s", split[3], split[1])
}
