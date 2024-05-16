package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationEndpointDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_notification_endpoint" "notification_endpoint" {
  name = "Notification Endpoint Unique"
  protocol = "slack"
  url = "https://x.y"
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
