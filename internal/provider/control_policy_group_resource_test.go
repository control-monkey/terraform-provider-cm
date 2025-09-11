package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfControlPolicyGroupResource     = "cm_control_policy_group"
	controlPolicyGroupTfResourceName = "control_policy_group"

	controlPolicyGroupName        = "tf control policy group"
	controlPolicyGroupDescription = "test"

	controlPolicySeverityAfterUpdate  = "high"
	controlPolicyGroupNameAfterUpdate = "updated tf control policy group"
)

func testAccControlPolicyGroupResourceSetup() string {
	return `
resource "cm_control_policy" "control_policy_test" {
  name        = "AWS Resources should have the Env tag with value Dev/Stage/Prod"
  description = "All AWS infrastructure should have the Env tag with value Dev/Stage/Prod."
  type        = "aws_required_tags"
  parameters  = jsonencode({
    tags = [
      {
        key           = "Env"
        allowedValues = [
          "Dev",
          "Stage",
          "Prod"
        ]
      }
    ]
  })
}
`
}

func TestAccControlPolicyGroupResource(t *testing.T) {
	// Test environment variables used by this function
	controlPolicyId := test_config.GetControlPolicyId()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccControlPolicyGroupResourceSetup() + providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	description = "%s"
	control_policies = [
		{
			control_policy_id = cm_control_policy.control_policy_test.id
		}
	]
}
`, tfControlPolicyGroupResource, controlPolicyGroupTfResourceName, controlPolicyGroupName, controlPolicyGroupDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "name", controlPolicyGroupName),
					resource.TestCheckResourceAttr(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "description", controlPolicyGroupDescription),
					resource.TestCheckResourceAttr(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "control_policies.#", "1"),
					resource.TestCheckResourceAttrPair(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "control_policies.0.control_policy_id", "cm_control_policy.control_policy_test", "id"),
					resource.TestCheckNoResourceAttr(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "control_policies.0.severity"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "id"),
				),
			},
			// Update and Read testing
			{
				ConfigVariables: config.Variables{
					//"control_policies": config.StringVariable(ControlPolicyGroupParametersAfterUpdate),
					"control_policies": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
						"control_policy_id": config.StringVariable(controlPolicyId),
						"severity":          config.StringVariable(controlPolicySeverityAfterUpdate),
					})),
				},
				Config: testAccControlPolicyGroupResourceSetup() + providerConfig + fmt.Sprintf(`
variable "control_policies" {
	type = list(object({
				control_policy_id = string
				severity = optional(string)
			}))
}

resource "%s" "%s" {
	name = "%s"
	description = "%s"
	control_policies = var.control_policies
}
`, tfControlPolicyGroupResource, controlPolicyGroupTfResourceName, controlPolicyGroupNameAfterUpdate, controlPolicyGroupDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "name", controlPolicyGroupNameAfterUpdate),
					resource.TestCheckResourceAttr(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "description", controlPolicyGroupDescription),
					resource.TestCheckResourceAttr(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "control_policies.#", "1"),
					resource.TestCheckResourceAttr(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "control_policies.0.control_policy_id", controlPolicyId),
					resource.TestCheckResourceAttr(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "control_policies.0.severity", controlPolicySeverityAfterUpdate),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "id"),
				),
			},
			{
				ResourceName:      controlPolicyGroupResourceName(controlPolicyGroupTfResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ConfigVariables: config.Variables{
					//"control_policies": config.StringVariable(ControlPolicyGroupParametersAfterUpdate),
					"control_policies": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
						"control_policy_id": config.StringVariable(controlPolicyId),
						"severity":          config.StringVariable(controlPolicySeverityAfterUpdate),
					})),
				},
			},
		},
	})
}

func controlPolicyGroupResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfControlPolicyGroupResource, s)
}
