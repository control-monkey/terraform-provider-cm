package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccControlPolicyGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_control_policy" "control_policy" {
  name = "Control Policy"
  type = "aws_denied_regions"
  parameters = jsonencode({
    regions = ["us-east-1"]
  })
}

resource "cm_control_policy_group" "control_policy_group" {
  name = "Control Policy Group Unique"
  control_policies = [
	{
	  control_policy_id = cm_control_policy.control_policy.id
	}
  ]
}

data "cm_control_policy_group" "control_policy_group" {
  name = cm_control_policy_group.control_policy_group.name
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_control_policy_group.control_policy_group", "id"),
					resource.TestCheckResourceAttr("data.cm_control_policy_group.control_policy_group", "name", "Control Policy Group Unique"),
				),
			},
		},
	})
}
