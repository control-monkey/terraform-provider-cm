package provider

import (
	"fmt"
	"testing"

	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomAbacConfigurationDataSource(t *testing.T) {
	// Test environment variables used by this function
	orgId := test_config.GetOrgId()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				ConfigVariables: config.Variables{
					"roles": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
						"org_id":   config.StringVariable(orgId),
						"org_role": config.StringVariable(cmTypes.RoleViewer),
					})),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "roles" {
  type = list(object({
	org_id = string
	org_role = string
	team_ids = optional(list(string))
    })
  )
}

resource "cm_custom_abac_configuration" "custom_abac_configuration" {
  custom_abac_id = "xxxx"
  name = "Custom ABAC Configuration"
  roles = var.roles

}

data "cm_custom_abac_configuration" "custom_abac_configuration_data" {
  name = cm_custom_abac_configuration.custom_abac_configuration.name
}

data "cm_custom_abac_configuration" "custom_abac_configuration_data2" {
  id = cm_custom_abac_configuration.custom_abac_configuration.id
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_custom_abac_configuration.custom_abac_configuration_data", "id"),
					resource.TestCheckResourceAttr("data.cm_custom_abac_configuration.custom_abac_configuration_data", "name", "Custom ABAC Configuration"),
					resource.TestCheckResourceAttrPair("data.cm_custom_abac_configuration.custom_abac_configuration_data", "id", "cm_custom_abac_configuration.custom_abac_configuration", "id"),
					resource.TestCheckResourceAttrSet("data.cm_custom_abac_configuration.custom_abac_configuration_data2", "id"),
					resource.TestCheckResourceAttr("data.cm_custom_abac_configuration.custom_abac_configuration_data2", "name", "Custom ABAC Configuration"),
					resource.TestCheckResourceAttrPair("data.cm_custom_abac_configuration.custom_abac_configuration_data2", "id", "cm_custom_abac_configuration.custom_abac_configuration", "id"),
				),
			},
			{
				ConfigVariables: config.Variables{
					"custom_abac_configuration_name": config.StringVariable("Custom ABAC Configuration"),
					"roles": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
						"org_id":   config.StringVariable(orgId),
						"org_role": config.StringVariable(cmTypes.RoleViewer),
					})),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "roles" {
  type = list(object({
	org_id = string
	org_role = string
	team_ids = optional(list(string))
    })
  )
}

resource "cm_custom_abac_configuration" "custom_abac_configuration" {
  custom_abac_id = "xxxx"
  name = "Custom ABAC Configuration"
  roles = var.roles
}

variable "custom_abac_configuration_name" {
	type = string
}

data "cm_custom_abac_configuration" "custom_abac_configuration_data" {
  name = var.custom_abac_configuration_name
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_custom_abac_configuration.custom_abac_configuration_data", "id"),
					resource.TestCheckResourceAttr("data.cm_custom_abac_configuration.custom_abac_configuration_data", "name", "Custom ABAC Configuration"),
					resource.TestCheckResourceAttrPair("data.cm_custom_abac_configuration.custom_abac_configuration_data", "id", "cm_custom_abac_configuration.custom_abac_configuration", "id"),
				),
			},
		},
	})
}
