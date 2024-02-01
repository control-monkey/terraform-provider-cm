package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmUsersTeam = "cm_team_users"

	teamUsersResourceName = "team_users"
	teamId                = "team-pgublox37u"
	userEmail             = "example@email.com"
)

func TestAccTeamUsersResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  team_id = "%s"
  users = [
    {
      email = "%s"
    },
  ]
}
`, cmUsersTeam, teamUsersResourceName, teamId, userEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(teamUsersResource(teamUsersResourceName), "id"),
					resource.TestCheckResourceAttrSet(teamUsersResource(teamUsersResourceName), "team_id"),
					resource.TestCheckResourceAttr(teamUsersResource(teamUsersResourceName), "users.0.email", userEmail),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
 team_id = "%s"
}
`, cmUsersTeam, teamUsersResourceName, teamId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(teamUsersResource(teamUsersResourceName), "id"),
					resource.TestCheckResourceAttr(teamUsersResource(teamUsersResourceName), "team_id", teamId),
					resource.TestCheckNoResourceAttr(teamUsersResource(teamUsersResourceName), "users"),
				),
			},
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmUsersTeam, teamUsersResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(teamUsersResource(teamUsersResourceName), "id"),
					resource.TestCheckResourceAttr(teamUsersResource(teamUsersResourceName), "id", teamId),
					resource.TestCheckResourceAttr(teamUsersResource(teamUsersResourceName), "team_id", teamId),
					resource.TestCheckNoResourceAttr(teamUsersResource(teamUsersResourceName), "users"),
				),
			},
		},
	})
}

func teamUsersResource(s string) string {
	return fmt.Sprintf("%s.%s", cmUsersTeam, s)
}
