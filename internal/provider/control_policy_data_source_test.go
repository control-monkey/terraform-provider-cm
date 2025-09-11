package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccControlPolicyDataSource(t *testing.T) {
	// Test environment variables used by this function
	controlPolicyId := test_config.GetControlPolicyId()
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
					"control_policy_id": config.StringVariable(controlPolicyId),
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
					resource.TestCheckResourceAttr("data.cm_control_policy.control_policy", "id", controlPolicyId),
				),
			},
		},
	})
}
