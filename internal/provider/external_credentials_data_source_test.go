package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExternalCredentialDataSource(t *testing.T) {
	id := test_config.GetExternalCredentialId()
	name := test_config.GetExternalCredentialName()
	vendor := test_config.GetExternalCredentialVendor()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
data "cm_external_credential" "test" {
name = "%s"
vendor = "%s"
}`,
					name, vendor),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_external_credential.test", "id"),
					resource.TestCheckResourceAttr("data.cm_external_credential.test", "name", name),
					resource.TestCheckResourceAttr("data.cm_external_credential.test", "vendor", vendor),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
data "cm_external_credential" "test" {
id = "%s"
vendor = "%s"
}`,
					id, vendor),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_external_credential.test", "name"),
					resource.TestCheckResourceAttr("data.cm_external_credential.test", "id", id),
					resource.TestCheckResourceAttr("data.cm_external_credential.test", "vendor", vendor),
				),
			},
		},
	})
}
