package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfCmOrganizationPermissions = "cm_organization_permissions"

	opResourceName = "organization_permissions"
)

// Two organization roles, to prove a single team can hold more than one.
func testAccOrganizationPermissionsResourceSetup() string {
	return `
resource "cm_team" "op_team" {
  name = "Organization Permissions Team"
}

resource "cm_custom_role" "op_role" {
  name = "Organization Permissions Role"
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

resource "cm_custom_role" "op_role2" {
  name = "Organization Permissions Role 2"
  type = "organizationRole"

  permissions = [
    {
      names = ["org:namespace:read"]
    },
  ]
}
`
}

func TestAccOrganizationPermissionsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Grant two organization roles to one team
			{
				Config: providerConfig + testAccOrganizationPermissionsResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  team_id = cm_team.op_team.id
  permissions = [
	{
  	  custom_role_id = cm_custom_role.op_role.id
	},
	{
  	  custom_role_id = cm_custom_role.op_role2.id
	},
  ]
}
`, tfCmOrganizationPermissions, opResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(organizationPermissionsResourceName(opResourceName), "team_id"),
					resource.TestCheckResourceAttr(organizationPermissionsResourceName(opResourceName), "permissions.#", "2"),

					resource.TestCheckResourceAttrSet(organizationPermissionsResourceName(opResourceName), "id"),
					resource.TestCheckResourceAttrPair(organizationPermissionsResourceName(opResourceName), "id", "cm_team.op_team", "id"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			// Revoke one of them
			{
				Config: providerConfig + testAccOrganizationPermissionsResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  team_id = cm_team.op_team.id
  permissions = [
	{
  	  custom_role_id = cm_custom_role.op_role.id
	},
  ]
}
`, tfCmOrganizationPermissions, opResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(organizationPermissionsResourceName(opResourceName), "permissions.#", "1"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			{
				ResourceName:      organizationPermissionsResourceName(opResourceName),
				ImportStateVerify: true,
				ImportState:       true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(organizationPermissionsResourceName(opResourceName), "team_id"),
					resource.TestCheckResourceAttr(organizationPermissionsResourceName(opResourceName), "permissions.#", "1"),

					resource.TestCheckResourceAttrSet(organizationPermissionsResourceName(opResourceName), "id"),
				),
			},
		},
	})
}

func organizationPermissionsResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfCmOrganizationPermissions, s)
}
