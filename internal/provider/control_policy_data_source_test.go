package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccControlPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
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
			{
				ConfigVariables: config.Variables{
					"control_policy_id": config.StringVariable(os.Getenv("CMP_ID")),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "control_policy_id" {
  type = string
}

data "cm_control_policy" "control_policy" {
  id = var.control_policy_id
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_control_policy.control_policy", "id"),
					resource.TestCheckResourceAttrSet("data.cm_control_policy.control_policy", "name"),
					resource.TestCheckResourceAttr("data.cm_control_policy.control_policy", "id", os.Getenv("CMP_ID")),
				),
			},
		},
	})
}
