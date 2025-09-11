package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomRoleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_custom_role" "custom_role" {
  name = "Custom Role"
}

data "cm_custom_role" "custom_role_data" {
  name = cm_custom_role.custom_role.name
}

data "cm_custom_role" "custom_role_data2" {
  id = cm_custom_role.custom_role.id
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_custom_role.custom_role_data", "id"),
					resource.TestCheckResourceAttr("data.cm_custom_role.custom_role_data", "name", "Custom Role"),
					resource.TestCheckResourceAttrPair("data.cm_custom_role.custom_role_data", "id", "cm_custom_role.custom_role", "id"),
					resource.TestCheckResourceAttrSet("data.cm_custom_role.custom_role_data2", "id"),
					resource.TestCheckResourceAttr("data.cm_custom_role.custom_role_data2", "name", "Custom Role"),
					resource.TestCheckResourceAttrPair("data.cm_custom_role.custom_role_data2", "id", "cm_custom_role.custom_role", "id"),
				),
			},
			{
				ConfigVariables: config.Variables{
					"role_name": config.StringVariable("Custom Role"),
				},
				Config: providerConfig + fmt.Sprintf(`
resource "cm_custom_role" "custom_role" {
  name = "Custom Role"
}

variable "role_name" {
	type = string
}

data "cm_custom_role" "custom_role_data" {
  name = var.role_name
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_custom_role.custom_role_data", "id"),
					resource.TestCheckResourceAttr("data.cm_custom_role.custom_role_data", "name", "Custom Role"),
					resource.TestCheckResourceAttrPair("data.cm_custom_role.custom_role_data", "id", "cm_custom_role.custom_role", "id"),
				),
			},
		},
	})
}
