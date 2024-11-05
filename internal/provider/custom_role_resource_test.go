package provider

import (
	"fmt"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfCustomRoleResource     = "cm_custom_role"
	customRoleTfResourceName = "custom_role"

	customRoleName        = "Create Stack"
	customRoleDescription = "test"

	customRolePermission1 = "stack:create"
	customRolePermission2 = "stack:createFromTemplate"

	customRoleStackRestriction = "restrictReadToOwnStacks"

	customRoleNameAfterUpdate = "updated name"
)

func TestAccCustomRoleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	description = "%s"
	permissions = [
		{
			name = "%s"
		},
	]
	stack_restriction = "%s"
}
`, tfCustomRoleResource, customRoleTfResourceName, customRoleName, customRoleDescription, customRolePermission1, customRoleStackRestriction),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "name", customRoleName),
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "description", customRoleDescription),
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "permissions.#", "1"),
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "permissions.0.name", customRolePermission1),

					resource.TestCheckResourceAttrSet(customRoleResourceName(customRoleTfResourceName), "stack_restriction"),
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "stack_restriction", customRoleStackRestriction),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(customRoleResourceName(customRoleTfResourceName), "id"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{ // Update and Read testing
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	description = "%s"
	permissions = [
		{
			name = "%s"
		},
		{
			name = "%s"
		}
	]
	stack_restriction = "%s"
}
`, tfCustomRoleResource, customRoleTfResourceName, customRoleNameAfterUpdate, customRoleDescription, customRolePermission1, customRolePermission2, customRoleStackRestriction),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "name", customRoleNameAfterUpdate),
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "description", customRoleDescription),
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "permissions.#", "2"),
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "permissions.0.name", customRolePermission1),
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "permissions.1.name", customRolePermission2),
					resource.TestCheckResourceAttrSet(customRoleResourceName(customRoleTfResourceName), "stack_restriction"),
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "stack_restriction", customRoleStackRestriction),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(customRoleResourceName(customRoleTfResourceName), "id"),

					//resource.TestCheckNoResourceAttr(customRoleResourceName(customRoleTfResourceName), "stack_restriction"),
				),
			},
			{
				ConfigVariables: config.Variables{
					"permission_name": config.StringVariable(customRolePermission2),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "permission_name" {
	type = string
}

resource "%s" "%s" {
	name = "%s"
	permissions = [
		{
			name = var.permission_name
		},
	]
}
`, tfCustomRoleResource, customRoleTfResourceName, customRoleNameAfterUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "name", customRoleNameAfterUpdate),
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "permissions.#", "1"),
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "permissions.0.name", customRolePermission2),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(customRoleResourceName(customRoleTfResourceName), "id"),

					resource.TestCheckNoResourceAttr(customRoleResourceName(customRoleTfResourceName), "description"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				ConfigVariables: config.Variables{
					"permission_name": config.StringVariable(customRolePermission2),
				},
				ResourceName:      customRoleResourceName(customRoleTfResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "name", customRoleNameAfterUpdate),
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "permissions.#", "1"),
					resource.TestCheckResourceAttr(customRoleResourceName(customRoleTfResourceName), "permissions.0.name", customRolePermission2),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(customRoleResourceName(customRoleTfResourceName), "id"),

					resource.TestCheckNoResourceAttr(customRoleResourceName(customRoleTfResourceName), "description"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
		},
	})
}

func customRoleResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfCustomRoleResource, s)
}
