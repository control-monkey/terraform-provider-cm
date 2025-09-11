package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmUsersTeam = "cm_team_users"

	teamUsersResourceName = "team_users"
	userEmail             = "example@email.com"
)

func testAccTeamUsersResourceSetup() string {
	return `
resource "cm_team" "test_team" {
  name = "TestTeam"
}
`
}

func TestAccTeamUsersResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccTeamUsersResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  team_id = cm_team.test_team.id
  users = [
    {
      email = "%s"
    },
  ]
}
`, cmUsersTeam, teamUsersResourceName, userEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(teamUsersResource(teamUsersResourceName), "id"),
					resource.TestCheckResourceAttrSet(teamUsersResource(teamUsersResourceName), "team_id"),
					resource.TestCheckResourceAttr(teamUsersResource(teamUsersResourceName), "users.0.email", userEmail),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccTeamUsersResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
 team_id = cm_team.test_team.id
}
`, cmUsersTeam, teamUsersResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(teamUsersResource(teamUsersResourceName), "id"),
					resource.TestCheckResourceAttrSet(teamUsersResource(teamUsersResourceName), "team_id"),
					resource.TestCheckNoResourceAttr(teamUsersResource(teamUsersResourceName), "users"),
				),
			},
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmUsersTeam, teamUsersResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(teamUsersResource(teamUsersResourceName), "id"),
					resource.TestCheckResourceAttrSet(teamUsersResource(teamUsersResourceName), "team_id"),
					resource.TestCheckNoResourceAttr(teamUsersResource(teamUsersResourceName), "users"),
				),
			},
		},
	})
}

func teamUsersResource(s string) string {
	return fmt.Sprintf("%s.%s", cmUsersTeam, s)
}
