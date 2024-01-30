package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmTeam = "cm_team"

	teamResourceName = "team"
	teamName         = "Dev Team"
	teamCustomIdpId  = "t123"

	teamNameAfterUpdate = "Prod team"
)

func TestAccTeamResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
 name = "%s"
 custom_idp_id = "%s"
}
`, cmTeam, teamResourceName, teamName, teamCustomIdpId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(teamResource(teamResourceName), "name", teamName),
					resource.TestCheckResourceAttr(teamResource(teamResourceName), "custom_idp_id", teamCustomIdpId),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(teamResource(teamResourceName), "id"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
 name = "%s"
}
`, cmTeam, teamResourceName, teamNameAfterUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(teamResource(teamResourceName), "id"),
					resource.TestCheckResourceAttr(teamResource(teamResourceName), "name", teamNameAfterUpdate),
					resource.TestCheckNoResourceAttr(teamResource(teamResourceName), "custom_idp_id"),
				),
			},
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmTeam, teamResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func teamResource(s string) string {
	return fmt.Sprintf("%s.%s", cmTeam, s)
}
