package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmNamespacePermissions = "cm_namespace_permissions"

	namespacePermissionsResourceName = "namespace_permissions"
	namespaceId                      = "ns-x82yjdyahc"
	permissionUsername               = "Registry Acceptance Test"
	permissionRoleViewer             = "viewer"
	permissionRoleAdmin              = "admin"
)

func TestAccNamespacePermissionsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  namespace_id = "%s"
  permissions = [
    {
      programmatic_username = "%s"
	  role = "%s"
    },
  ]
}
`, cmNamespacePermissions, namespacePermissionsResourceName, namespaceId, permissionUsername, permissionRoleViewer),
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
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
 namespace_id = "%s"

  permissions = [
    {
      programmatic_username = "%s"
	  role = "%s"
    },
  ]
}
`, cmNamespacePermissions, namespacePermissionsResourceName, namespaceId, permissionUsername, permissionRoleAdmin),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(namespacePermissionsResource(namespacePermissionsResourceName), "id"),
					resource.TestCheckResourceAttr(namespacePermissionsResource(namespacePermissionsResourceName), "namespace_id", namespaceId),
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
					resource.TestCheckResourceAttr(namespacePermissionsResource(namespacePermissionsResourceName), "id", namespaceId),
					resource.TestCheckResourceAttr(namespacePermissionsResource(namespacePermissionsResourceName), "namespace_id", namespaceId),
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
