package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccControlPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_control_policy" "control_policy" {
  name = "Control Policy Unique"
  type = "aws_denied_regions"
  parameters = jsonencode({
    regions = ["us-east-1"]
  })
}

data "cm_control_policy" "control_policy" {
  name = cm_control_policy.control_policy.name
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_control_policy.control_policy", "id"),
					resource.TestCheckResourceAttr("data.cm_control_policy.control_policy", "name", "Control Policy Unique"),
				),
			},
		},
	})
}
