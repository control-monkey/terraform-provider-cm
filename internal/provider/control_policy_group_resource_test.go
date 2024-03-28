package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"testing"

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

func TestAccControlPolicyGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	description = "%s"
	control_policies = [
		{
			control_policy_id = "%s"
		}
	]
}
`, tfControlPolicyGroupResource, controlPolicyGroupTfResourceName, controlPolicyGroupName, controlPolicyGroupDescription, controlPolicyId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "name", controlPolicyGroupName),
					resource.TestCheckResourceAttr(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "description", controlPolicyGroupDescription),
					resource.TestCheckResourceAttr(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "control_policies.#", "1"),
					resource.TestCheckResourceAttr(controlPolicyGroupResourceName(controlPolicyGroupTfResourceName), "control_policies.0.control_policy_id", controlPolicyId),
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
				Config: providerConfig + fmt.Sprintf(`
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
