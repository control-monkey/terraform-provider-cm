package provider

import (
	"fmt"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmNotificationSlackApp  = "cm_notification_slack_app"
	slackAppTfResourceName  = "slack_app"
	slackAppName            = "tf-acc-slack-app"
	slackAppNameAfterUpdate = "tf-acc-slack-app-updated"
)

func TestAccNotificationSlackAppResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name = "%s"
  bot_auth_token = "xoxb-***"
}
`, cmNotificationSlackApp, slackAppTfResourceName, slackAppName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(slackAppResourceName(slackAppTfResourceName), "id"),
					resource.TestCheckResourceAttr(slackAppResourceName(slackAppTfResourceName), "name", slackAppName),
					resource.TestCheckResourceAttr(slackAppResourceName(slackAppTfResourceName), "bot_auth_token", "xoxb-***"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name = "%s"
  bot_auth_token = "ignored-token"

  lifecycle {
    ignore_changes = [bot_auth_token]
  }

}
`, cmNotificationSlackApp, slackAppTfResourceName, slackAppNameAfterUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(slackAppResourceName(slackAppTfResourceName), "id"),
					resource.TestCheckResourceAttr(slackAppResourceName(slackAppTfResourceName), "name", slackAppNameAfterUpdate),
					resource.TestCheckResourceAttr(slackAppResourceName(slackAppTfResourceName), "bot_auth_token", "xoxb-***"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name = "%s"
  bot_auth_token = "xoxb-*****"
}
`, cmNotificationSlackApp, slackAppTfResourceName, slackAppNameAfterUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(slackAppResourceName(slackAppTfResourceName), "id"),
					resource.TestCheckResourceAttr(slackAppResourceName(slackAppTfResourceName), "name", slackAppNameAfterUpdate),
					resource.TestCheckResourceAttr(slackAppResourceName(slackAppTfResourceName), "bot_auth_token", "xoxb-*****"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				ResourceName: fmt.Sprintf("%s.%s", cmNotificationSlackApp, slackAppTfResourceName),
				ImportState:  true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(slackAppResourceName(slackAppTfResourceName), "id"),
					resource.TestCheckResourceAttr(slackAppResourceName(slackAppTfResourceName), "name", slackAppNameAfterUpdate),
					resource.TestCheckNoResourceAttr(slackAppResourceName(slackAppTfResourceName), "bot_auth_token"),
				),
			},
		},
	})
}

func slackAppResourceName(s string) string {
	return fmt.Sprintf("%s.%s", cmNotificationSlackApp, s)
}
