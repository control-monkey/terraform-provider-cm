package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"
	"github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfCustomAbacConfigurationResourceResource = "cm_custom_abac_configuration"
	customAbacConfigurationTfResourceName     = "custom_abac_configuration"

	customAbacConfigurationAbacId = "abc-123"
	customAbacConfigurationName   = "ABAC"

	customAbacConfigurationOrgRole  = "member"
	customAbacConfigurationOrgRole2 = "viewer"

	customAbacConfigurationNameAfterUpdate = "updated name"
)

func TestAccCustomAbacConfigurationResourceResource(t *testing.T) {
	// Test environment variables used by this function
	orgId := test_config.GetOrgId()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_team" "team1" {
	name = "Developers"
}

resource "cm_team" "team2" {
	name = "QA"
}

resource "%s" "%s" {
	custom_abac_id = "%s"
	name = "%s"
	roles = [
		{
			org_id = "%s"
			org_role = "%s"
			team_ids = [cm_team.team1.id, cm_team.team2.id]
		}
	]
}
`, tfCustomAbacConfigurationResourceResource, customAbacConfigurationTfResourceName, customAbacConfigurationAbacId, customAbacConfigurationName,
					orgId, customAbacConfigurationOrgRole),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "custom_abac_id", customAbacConfigurationAbacId),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "name", customAbacConfigurationName),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.#", "1"),
					resource.TestCheckResourceAttrSet(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.0.org_id"),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.0.org_role", customAbacConfigurationOrgRole),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.0.team_ids.#", "2"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_team" "team1" {
	name = "Developers"
}

resource "cm_team" "team2" {
	name = "QA"
}

resource "%s" "%s" {
	custom_abac_id = "%s"
	name = "%s"
	roles = [
		{
			org_id = "%s"
			org_role = "%s"
			team_ids = [cm_team.team1.id, cm_team.team2.id]
		}
	]
}
`, tfCustomAbacConfigurationResourceResource, customAbacConfigurationTfResourceName, customAbacConfigurationAbacId, customAbacConfigurationName,
					orgId, customAbacConfigurationOrgRole2),
				ExpectError: regexp.MustCompile(validationError),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_team" "team1" {
	name = "Developers"
}

resource "cm_team" "team2" {
	name = "QA"
}

resource "%s" "%s" {
	custom_abac_id = "%s"
	name = "%s"
	roles = [
		{
			org_id = "%s"
			org_role = "%s"
			team_ids = [cm_team.team1.id]
		}
	]
}
`, tfCustomAbacConfigurationResourceResource, customAbacConfigurationTfResourceName, customAbacConfigurationAbacId, customAbacConfigurationNameAfterUpdate,
					orgId, customAbacConfigurationOrgRole),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "custom_abac_id", customAbacConfigurationAbacId),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "name", customAbacConfigurationNameAfterUpdate),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.#", "1"),
					resource.TestCheckResourceAttrSet(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.0.org_id"),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.0.org_role", customAbacConfigurationOrgRole),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.0.team_ids.#", "1"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				ConfigVariables: config.Variables{
					"org_id": config.StringVariable(orgId),
				},
				Config: providerConfig + fmt.Sprintf(`
resource "cm_team" "team1" {
	name = "Developers"
}

resource "cm_team" "team2" {
	name = "QA"
}

variable "org_id" {
	type = string
}

resource "%s" "%s" {
	custom_abac_id = "%s"
	name = "%s"
	roles = [
		{
			org_id = var.org_id
			org_role = "%s"
		}
	]
}
`, tfCustomAbacConfigurationResourceResource, customAbacConfigurationTfResourceName, customAbacConfigurationAbacId, customAbacConfigurationNameAfterUpdate, customAbacConfigurationOrgRole2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "custom_abac_id", customAbacConfigurationAbacId),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "name", customAbacConfigurationNameAfterUpdate),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.#", "1"),
					resource.TestCheckResourceAttrSet(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.0.org_id"),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.0.org_role", customAbacConfigurationOrgRole2),
					resource.TestCheckNoResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.0.team_ids"),
				),
			},
			test_helpers.GetValidateNoDriftStep(),
			{
				ConfigVariables: config.Variables{
					"org_id": config.StringVariable(orgId),
				},
				ResourceName:      customAbacConfigurationResourceName(customAbacConfigurationTfResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "custom_abac_id", customAbacConfigurationAbacId),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "name", customAbacConfigurationNameAfterUpdate),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.#", "1"),
					resource.TestCheckResourceAttrSet(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.0.org_id"),
					resource.TestCheckResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.0.org_role", customAbacConfigurationOrgRole2),
					resource.TestCheckNoResourceAttr(customAbacConfigurationResourceName(customAbacConfigurationTfResourceName), "roles.0.team_ids"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
		},
	})
}

func customAbacConfigurationResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfCustomAbacConfigurationResourceResource, s)
}
