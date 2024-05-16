package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTeamDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_team" "team" {
  name = "Team Unique"
}

data "cm_team" "team" {
  name = cm_team.team.name
}`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_team.team", "id"),
					resource.TestCheckResourceAttr("data.cm_team.team", "name", "Team Unique"),
				),
			},
		},
	})
}
