package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmNamespacePermissions = "cm_namespace_permissions"

	namespacePermissionsResourceName = "namespace_permissions"
	permissionUsername               = "Registry Acceptance Test"
	permissionRoleViewer             = "viewer"
	permissionRoleAdmin              = "admin"
)

func testAccNamespacePermissionsResourceSetup() string {
	return `
resource "cm_namespace" "test_namespace" {
  name = "TestNamespace"
}
`
}

func TestAccNamespacePermissionsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccNamespacePermissionsResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  namespace_id = cm_namespace.test_namespace.id
  permissions = [
    {
      programmatic_username = "%s"
	  role = "%s"
    },
  ]
}
`, cmNamespacePermissions, namespacePermissionsResourceName, permissionUsername, permissionRoleViewer),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(namespacePermissionsResource(namespacePermissionsResourceName), "id"),
					resource.TestCheckResourceAttrSet(namespacePermissionsResource(namespacePermissionsResourceName), "namespace_id"),
					resource.TestCheckResourceAttr(namespacePermissionsResource(namespacePermissionsResourceName), "permissions.0.programmatic_username", permissionUsername),
					resource.TestCheckResourceAttr(namespacePermissionsResource(namespacePermissionsResourceName), "permissions.0.role", permissionRoleViewer),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccNamespacePermissionsResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
 namespace_id = cm_namespace.test_namespace.id

  permissions = [
    {
      programmatic_username = "%s"
	  role = "%s"
    },
  ]
}
`, cmNamespacePermissions, namespacePermissionsResourceName, permissionUsername, permissionRoleAdmin),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(namespacePermissionsResource(namespacePermissionsResourceName), "id"),
					resource.TestCheckResourceAttrSet(namespacePermissionsResource(namespacePermissionsResourceName), "namespace_id"),
					resource.TestCheckResourceAttr(namespacePermissionsResource(namespacePermissionsResourceName), "permissions.0.programmatic_username", permissionUsername),
					resource.TestCheckResourceAttr(namespacePermissionsResource(namespacePermissionsResourceName), "permissions.0.role", permissionRoleAdmin),
				),
			},
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmNamespacePermissions, namespacePermissionsResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(namespacePermissionsResource(namespacePermissionsResourceName), "id"),
					resource.TestCheckResourceAttrSet(namespacePermissionsResource(namespacePermissionsResourceName), "namespace_id"),
					resource.TestCheckResourceAttr(namespacePermissionsResource(namespacePermissionsResourceName), "permissions.0.programmatic_username", permissionUsername),
					resource.TestCheckResourceAttr(namespacePermissionsResource(namespacePermissionsResourceName), "permissions.0.role", permissionRoleAdmin),
				),
			},
		},
	})
}

func namespacePermissionsResource(s string) string {
	return fmt.Sprintf("%s.%s", cmNamespacePermissions, s)
}
