package organizationPermissions

import (
	"fmt"
	"strings"

	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID          types.String       `tfsdk:"id"`
	TeamId      types.String       `tfsdk:"team_id"`
	Permissions []*PermissionModel `tfsdk:"permissions"`
}

type PermissionModel struct {
	CustomRoleId types.String `tfsdk:"custom_role_id"`
}

// Hash covers every field: there are no mutable fields, so every change is a delete + create.
func (e *PermissionModel) Hash() string {
	retVal := ""

	if e.CustomRoleId.IsNull() == false {
		retVal += fmt.Sprintf("CustomRoleId:%s:", e.CustomRoleId.ValueString())
	}

	return retVal
}

func (e *PermissionModel) GetBlockIdentifier() string {
	retVal := ""

	if helpers.IsKnown(e.CustomRoleId) {
		retVal = fmt.Sprintf("CustomRoleId:%s", e.CustomRoleId.ValueString())
	}

	return retVal
}

func CleanIdentifier(s string) string {
	split := strings.Split(s, ":")
	return split[1]
}
