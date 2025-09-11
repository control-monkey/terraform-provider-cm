package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationEndpointDataSource(t *testing.T) {
	// Test environment variables used by this function
	slackWebhookUrl := test_config.GetSlackWebhookUrl()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Read testing
			{
				ConfigVariables: config.Variables{
					"slack_url": config.StringVariable(slackWebhookUrl),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "slack_url" {
  type = string
}

resource "cm_notification_endpoint" "notification_endpoint" {
  name = "Notification Endpoint Unique"
  protocol = "slack"
  url = var.slack_url
}

data "cm_notification_endpoint" "notification_endpoint" {
  name = cm_notification_endpoint.notification_endpoint.name
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_notification_endpoint.notification_endpoint", "id"),
					resource.TestCheckResourceAttr("data.cm_notification_endpoint.notification_endpoint", "name", "Notification Endpoint Unique"),
				),
			},
		},
	})
}
