package namespace_permissions

import (
	"fmt"
	"github.com/control-monkey/controlmonkey-sdk-go/services/namespace_permissions"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID          types.String        `tfsdk:"id"`
	NamespaceId types.String        `tfsdk:"namespace_id"`
	Permissions []*PermissionsModel `tfsdk:"permissions"`
}

type PermissionsModel struct { //When new field is added consider Hash() function
	UserEmail            types.String `tfsdk:"user_email"`
	ProgrammaticUserName types.String `tfsdk:"programmatic_username"`
	TeamId               types.String `tfsdk:"team_id"`
	Role                 types.String `tfsdk:"role"`
	CustomRoleId         types.String `tfsdk:"custom_role_id"`
}

type MergedEntities struct {
	EntitiesToCreate []*namespace_permissions.NamespacePermission
	EntitiesToUpdate []*namespace_permissions.NamespacePermission
	EntitiesToDelete []*namespace_permissions.NamespacePermission
}

func (e *PermissionsModel) Hash() string {
	retVal := ""

	if e.UserEmail.IsNull() == false {
		retVal += fmt.Sprintf("UserEmail:%s:", e.UserEmail.ValueString())
	}
	if e.ProgrammaticUserName.IsNull() == false {
		retVal += fmt.Sprintf("ProgrammaticUserName:%s:", e.ProgrammaticUserName.ValueString())
	}
	if e.TeamId.IsNull() == false {
		retVal += fmt.Sprintf("TeamId:%s:", e.TeamId.ValueString())
	}
	if e.Role.IsNull() == false {
		retVal += fmt.Sprintf("Role:%s:", e.Role.ValueString())
	}
	if e.CustomRoleId.IsNull() == false {
		retVal += fmt.Sprintf("CustomRoleId:%s:", e.CustomRoleId.ValueString())
	}

	return retVal
}

func (e *PermissionsModel) GetBlockIdentifier() string {
	retVal := ""

	if e.UserEmail.IsNull() == false {
		retVal += fmt.Sprintf("UserEmail:%s:", e.UserEmail.ValueString())
	}
	if e.ProgrammaticUserName.IsNull() == false {
		retVal += fmt.Sprintf("ProgrammaticUserName:%s:", e.ProgrammaticUserName.ValueString())
	}
	if e.TeamId.IsNull() == false {
		retVal += fmt.Sprintf("TeamId:%s:", e.TeamId.ValueString())
	}

	return retVal
}
