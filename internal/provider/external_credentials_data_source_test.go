package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExternalCredentialsDataSource(t *testing.T) {
	id := test_config.GetExternalCredentialsId()
	name := test_config.GetExternalCredentialsName()
	vendor := test_config.GetExternalCredentialsVendor()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
data "cm_external_credentials" "test" {
name = "%s"
vendor = "%s"
}`,
					name, vendor),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_external_credentials.test", "id"),
					resource.TestCheckResourceAttr("data.cm_external_credentials.test", "name", name),
					resource.TestCheckResourceAttr("data.cm_external_credentials.test", "vendor", vendor),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
data "cm_external_credentials" "test" {
id = "%s"
vendor = "%s"
}`,
					id, vendor),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_external_credentials.test", "name"),
					resource.TestCheckResourceAttr("data.cm_external_credentials.test", "id", id),
					resource.TestCheckResourceAttr("data.cm_external_credentials.test", "vendor", vendor),
				),
			},
		},
	})
}
