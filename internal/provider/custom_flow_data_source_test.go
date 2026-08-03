package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	dataCmCustomFlow = "data.cm_custom_flow"

	cfDataResourceName = "custom_flow_by_name"
	cfDataName         = "tf custom flow data source"
)

func TestAccCustomFlowDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Look the custom flow up by name and confirm the id matches the managed resource
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name       = "%s"
  is_enabled = true
  flow_yaml  = "version: 1"
}

data "cm_custom_flow" "%s" {
  name       = "%s"
  depends_on = [%s.%s]
}
`, tfCmCustomFlow, cfDataResourceName, cfDataName, cfDataResourceName, cfDataName, tfCmCustomFlow, cfDataResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(customFlowDataSourceName(cfDataResourceName), "id"),
					resource.TestCheckResourceAttr(customFlowDataSourceName(cfDataResourceName), "name", cfDataName),
					resource.TestCheckResourceAttrPair(customFlowDataSourceName(cfDataResourceName), "id",
						customFlowResourceName(cfDataResourceName), "id"),
				),
			},
		},
	})
}

func customFlowDataSourceName(s string) string {
	return fmt.Sprintf("%s.%s", dataCmCustomFlow, s)
}
