package provider

import (
	"fmt"
	"testing"

	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmNotificationEndpoint = "cm_notification_endpoint"

	notificationEndpointResourceName = "notificationEndpoint"
	notificationEndpointName         = "Dev Endpoint"
	notificationEndpointProtocol     = cmTypes.SlackProtocol
	notificationEndpointUrl          = "https://hooks.slack.com"

	notificationEndpointNameAfterUpdate = "Prod notificationEndpoint"
)

func TestAccNotificationEndpointResource(t *testing.T) {
	// Test environment variables used by this function
	email1 := test_config.GetNotificationEndpointEmail1()
	email2 := test_config.GetNotificationEndpointEmail2()
	slackAppId := test_config.GetSlackAppId()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Step 1: Create Slack App endpoint (assumes env vars are set)
			{
				ConfigVariables: config.Variables{
					"slack_config": config.ObjectVariable(map[string]config.Variable{
						"notification_slack_app_id": config.StringVariable(slackAppId),
						"channel_id":                config.StringVariable("C123"),
					}),
					"emails": config.ListVariable(
						config.StringVariable(email1),
						config.StringVariable(email2),
					),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "slack_config" {
  type = object({
	notification_slack_app_id = string
	channel_id = string
  })
}

resource "%s" "%s" {
  name     = "%s"
  protocol = "%s"
  slack_app_config = var.slack_config
}
`, cmNotificationEndpoint, notificationEndpointResourceName, notificationEndpointName, cmTypes.SlackAppProtocol),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "name", notificationEndpointName),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "protocol", cmTypes.SlackAppProtocol),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "slack_app_config.notification_slack_app_id", slackAppId),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "slack_app_config.channel_id", "C123"),
					resource.TestCheckResourceAttrSet(notificationEndpointResource(notificationEndpointResourceName), "id"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			// Step 2: Update name (still Slack App)
			{
				ConfigVariables: config.Variables{
					"slack_app_id": config.StringVariable(slackAppId),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "slack_app_id" {
  type = string 
}

resource "%s" "%s" {
  name     = "%s"
  protocol = "%s"
  slack_app_config = {
    notification_slack_app_id = var.slack_app_id
    channel_id                = "C123"
  }
}
`, cmNotificationEndpoint, notificationEndpointResourceName, notificationEndpointNameAfterUpdate, cmTypes.SlackAppProtocol),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(notificationEndpointResource(notificationEndpointResourceName), "id"),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "name", notificationEndpointNameAfterUpdate),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "protocol", cmTypes.SlackAppProtocol),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "slack_app_config.notification_slack_app_id", slackAppId),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "slack_app_config.channel_id", "C123"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			// Step 3: Switch to Email protocol and set email_addresses
			{
				ConfigVariables: config.Variables{
					"email1": config.StringVariable(email1),
					"email2": config.StringVariable(email2),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "email1" {
  type = string 
}

variable "email2" {
  type = string 
}

resource "%s" "%s" {
  name            = "%s"
  protocol        = "%s"
  email_addresses = [var.email1, var.email2]
}
`, cmNotificationEndpoint, notificationEndpointResourceName, notificationEndpointNameAfterUpdate, cmTypes.EmailProtocol),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "protocol", cmTypes.EmailProtocol),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "email_addresses.#", "2"),
					resource.TestCheckNoResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "url"),
					resource.TestCheckNoResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "slack_app_config"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			// Step 4: Switch to Slack webhook (url)
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name     = "%s"
  protocol = "%s"
  url      = "%s"
}
`, cmNotificationEndpoint, notificationEndpointResourceName, notificationEndpointNameAfterUpdate, notificationEndpointProtocol, notificationEndpointUrl),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "protocol", notificationEndpointProtocol),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "url", notificationEndpointUrl),
					resource.TestCheckNoResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "slack_app_config"),
					resource.TestCheckNoResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "email_addresses"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			// Step 5: Import
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmNotificationEndpoint, notificationEndpointResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(notificationEndpointResource(notificationEndpointResourceName), "id"),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "name", notificationEndpointNameAfterUpdate),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "protocol", notificationEndpointProtocol),
				),
			},
		},
	})
}

func notificationEndpointResource(s string) string {
	return fmt.Sprintf("%s.%s", cmNotificationEndpoint, s)
}
