package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"
	"github.com/hashicorp/terraform-plugin-testing/config"

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

const (
	orgRoleTfResourceName = "organization_custom_role"
	orgRoleName           = "tf organization role"
)

// TestAccCustomRoleOrganizationRoleResource covers type, permissions.names and
// permissions.restrictions. Restrictions use a cloudAccount action - the API rejects them on org:*
// actions - and org:stack:create requires org:stack:read in the same block.
func TestAccCustomRoleOrganizationRoleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	type = "organizationRole"
	permissions = [
		{
			names = ["org:stack:read", "org:namespace:read", "org:stack:create"]
		},
		{
			names = ["cloudAccount:insights"]
			restrictions = [
				{
					cloud_provider   = "aws"
					cloud_account_id = "123456789012"
				}
			]
		},
	]
}
`, tfCustomRoleResource, orgRoleTfResourceName, orgRoleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(customRoleResourceName(orgRoleTfResourceName), "id"),
					resource.TestCheckResourceAttr(customRoleResourceName(orgRoleTfResourceName), "type", "organizationRole"),
					resource.TestCheckResourceAttr(customRoleResourceName(orgRoleTfResourceName), "permissions.#", "2"),
					// Order must round-trip exactly, in the order written above
					resource.TestCheckResourceAttr(customRoleResourceName(orgRoleTfResourceName), "permissions.0.names.#", "3"),
					resource.TestCheckResourceAttr(customRoleResourceName(orgRoleTfResourceName), "permissions.0.names.0", "org:stack:read"),
					resource.TestCheckResourceAttr(customRoleResourceName(orgRoleTfResourceName), "permissions.0.names.1", "org:namespace:read"),
					resource.TestCheckResourceAttr(customRoleResourceName(orgRoleTfResourceName), "permissions.0.names.2", "org:stack:create"),
					resource.TestCheckResourceAttr(customRoleResourceName(orgRoleTfResourceName), "permissions.1.restrictions.#", "1"),
					resource.TestCheckResourceAttr(customRoleResourceName(orgRoleTfResourceName), "permissions.1.restrictions.0.cloud_provider", "aws"),
					resource.TestCheckResourceAttr(customRoleResourceName(orgRoleTfResourceName), "permissions.1.restrictions.0.cloud_account_id", "123456789012"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			// Drop the restrictions block
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	type = "organizationRole"
	permissions = [
		{
			names = ["org:stack:read", "org:namespace:read", "org:stack:create"]
		},
	]
}
`, tfCustomRoleResource, orgRoleTfResourceName, orgRoleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(customRoleResourceName(orgRoleTfResourceName), "permissions.#", "1"),
					resource.TestCheckNoResourceAttr(customRoleResourceName(orgRoleTfResourceName), "permissions.0.restrictions"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			{
				ResourceName:      customRoleResourceName(orgRoleTfResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(customRoleResourceName(orgRoleTfResourceName), "id"),
					resource.TestCheckResourceAttr(customRoleResourceName(orgRoleTfResourceName), "type", "organizationRole"),
				),
			},
		},
	})
}
