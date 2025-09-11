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
resource "cm_notification_slack_app" "slack_app" {
  name = "Unique Slack App Name 123"
  bot_auth_token = "xoxb-***"
}
`),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_notification_slack_app" "slack_app" {
  name = "Unique Slack App Name 123"
  bot_auth_token = "xoxb-***"
}

data "%s" "%s" {
  name = cm_notification_slack_app.slack_app.name
}
`, cmNotificationSlackAppDataSource, slackAppTfDataSourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(slackAppDataSourceName(slackAppTfDataSourceName), "id"),
					resource.TestCheckResourceAttrSet(slackAppDataSourceName(slackAppTfDataSourceName), "name"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_notification_slack_app" "slack_app" {
  name = "Unique Slack App Name 123"
  bot_auth_token = "xoxb-***"
}

data "%s" "%s" {
  id = cm_notification_slack_app.slack_app.id
}
`, cmNotificationSlackAppDataSource, slackAppTfDataSourceName),
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
