package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmNotificationSlackAppDataSource = "cm_notification_slack_app"
	slackAppTfDataSourceName         = "slack_app"
)

func TestAccNotificationSlackAppDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name = "slack-app-data"
  bot_auth_token = "xoxb-***"
}
`, cmNotificationSlackAppDataSource, slackAppTfDataSourceName),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name = "slack-app-data"
  bot_auth_token = "xoxb-***"
}

data "%s" "%s" {
  name = "slack-app-data"
}
`, cmNotificationSlackAppDataSource, slackAppTfDataSourceName, cmNotificationSlackAppDataSource, slackAppTfDataSourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(slackAppDataSourceName(slackAppTfDataSourceName), "id"),
					resource.TestCheckResourceAttrSet(slackAppDataSourceName(slackAppTfDataSourceName), "name"),
				),
			},
		},
	})
}

func slackAppDataSourceName(s string) string {
	return fmt.Sprintf("%s.%s", cmNotificationSlackAppDataSource, s)
}
