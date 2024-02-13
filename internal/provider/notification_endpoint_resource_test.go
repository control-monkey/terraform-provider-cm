package provider

import (
	"fmt"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmNotificationEndpoint = "cm_notification_endpoint"

	notificationEndpointResourceName = "notificationEndpoint"
	notificationEndpointName         = "Dev Endpoint"
	notificationEndpointProtocol     = cmTypes.SlackProtocol
	notificationEndpointUrl          = "https://cm.cm"

	notificationEndpointNameAfterUpdate = "Prod notificationEndpoint"
)

func TestAccNotificationEndpointResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
 name = "%s"
 protocol = "%s"
 url = "%s"
}
`, cmNotificationEndpoint, notificationEndpointResourceName, notificationEndpointName, notificationEndpointProtocol, notificationEndpointUrl),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "name", notificationEndpointName),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "protocol", notificationEndpointProtocol),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "url", notificationEndpointUrl),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(notificationEndpointResource(notificationEndpointResourceName), "id"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
 name = "%s"
 protocol = "%s"
 url = "%s"
}
`, cmNotificationEndpoint, notificationEndpointResourceName, notificationEndpointNameAfterUpdate, notificationEndpointProtocol, notificationEndpointUrl),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(notificationEndpointResource(notificationEndpointResourceName), "id"),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "name", notificationEndpointNameAfterUpdate),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "protocol", notificationEndpointProtocol),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "url", notificationEndpointUrl),
				),
			},
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmNotificationEndpoint, notificationEndpointResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(notificationEndpointResource(notificationEndpointResourceName), "id"),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "name", notificationEndpointNameAfterUpdate),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "protocol", notificationEndpointProtocol),
					resource.TestCheckResourceAttr(notificationEndpointResource(notificationEndpointResourceName), "url", notificationEndpointUrl),
				),
			},
		},
	})
}

func notificationEndpointResource(s string) string {
	return fmt.Sprintf("%s.%s", cmNotificationEndpoint, s)
}
