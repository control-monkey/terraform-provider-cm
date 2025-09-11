package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmEventsSubscriptions = "cm_events_subscriptions"

	eventsSubscriptionsResourceName = "events_subscriptions"
	eventsSubscriptionsScope        = "namespace"
)

func testAccEventsSubscriptionsResourceSetup() string {
	// Test environment variables used by this function
	slackWebhookUrl := test_config.GetSlackWebhookUrl()

	return fmt.Sprintf(`
resource "cm_namespace" "namespace" {
  name = "namespace"
}

resource "cm_notification_endpoint" "notification_endpoint" {
  name = "test"
  protocol = "slack"
  url = "%s"
}

resource "cm_notification_endpoint" "notification_endpoint_2" {
  name = "test2"
  protocol = "slack"
  url = "%s"
}
`, slackWebhookUrl, slackWebhookUrl)
}

func TestAccEventsSubscriptionsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccEventsSubscriptionsResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  scope = "%s"
  scope_id = cm_namespace.namespace.id
  subscriptions = [
    {
      event_type = "stack::deployment::applyStarted"
	  notification_endpoint_id = cm_notification_endpoint.notification_endpoint.id
    },
  ]
}
`, cmEventsSubscriptions, eventsSubscriptionsResourceName, eventsSubscriptionsScope),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(eventsSubscriptionsResource(eventsSubscriptionsResourceName), "id"),
					resource.TestCheckResourceAttr(eventsSubscriptionsResource(eventsSubscriptionsResourceName), "scope", eventsSubscriptionsScope),
					resource.TestCheckResourceAttrSet(eventsSubscriptionsResource(eventsSubscriptionsResourceName), "scope_id"),
					resource.TestCheckResourceAttr(eventsSubscriptionsResource(eventsSubscriptionsResourceName), "subscriptions.#", "1"),
					resource.TestCheckResourceAttr(eventsSubscriptionsResource(eventsSubscriptionsResourceName), "subscriptions.0.event_type", "stack::deployment::applyStarted"),
					resource.TestCheckResourceAttrSet(eventsSubscriptionsResource(eventsSubscriptionsResourceName), "subscriptions.0.notification_endpoint_id"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccEventsSubscriptionsResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  scope = "%s"
  scope_id = cm_namespace.namespace.id

  subscriptions = [
    {
      event_type = "stack::deployment::applyFailure"
	  notification_endpoint_id = cm_notification_endpoint.notification_endpoint.id
    },
    {
      event_type = "stack::deployment::approvalTimeout"
	  notification_endpoint_id = cm_notification_endpoint.notification_endpoint.id
    },
    {
      event_type = "stack::deployment::applyStarted"
	  notification_endpoint_id = cm_notification_endpoint.notification_endpoint_2.id
    },
  ]
}
`, cmEventsSubscriptions, eventsSubscriptionsResourceName, eventsSubscriptionsScope),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(eventsSubscriptionsResource(eventsSubscriptionsResourceName), "id"),
					resource.TestCheckResourceAttr(eventsSubscriptionsResource(eventsSubscriptionsResourceName), "scope", eventsSubscriptionsScope),
					resource.TestCheckResourceAttrSet(eventsSubscriptionsResource(eventsSubscriptionsResourceName), "scope_id"),
					resource.TestCheckResourceAttr(eventsSubscriptionsResource(eventsSubscriptionsResourceName), "subscriptions.#", "3"),
				),
			},
			{
				ResourceName: fmt.Sprintf("%s.%s", cmEventsSubscriptions, eventsSubscriptionsResourceName),

				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(eventsSubscriptionsResource(eventsSubscriptionsResourceName), "id"),
					resource.TestCheckResourceAttr(eventsSubscriptionsResource(eventsSubscriptionsResourceName), "scope", eventsSubscriptionsScope),
					resource.TestCheckResourceAttrSet(eventsSubscriptionsResource(eventsSubscriptionsResourceName), "scope_id"),
					resource.TestCheckResourceAttr(eventsSubscriptionsResource(eventsSubscriptionsResourceName), "subscriptions.#", "3"),
				),
			},
		},
	})
}

func eventsSubscriptionsResource(s string) string {
	return fmt.Sprintf("%s.%s", cmEventsSubscriptions, s)
}
