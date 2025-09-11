package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfControlPolicyResource   = "cm_control_policy"
	ControlPolicyResourceName = "control_policy"

	ControlPolicyName        = "tf control policy"
	ControlPolicyDescription = "test"
	ControlPolicyType        = "aws_allowed_regions"

	ControlPolicyNameAfterUpdate       = "updated tf control policy"
	ControlPolicyParametersAfterUpdate = "{\"regions\":[\"us-east-1\"]}"
)

func TestAccControlPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	description = "%s"
	type = "%s"
	parameters = jsonencode({
	  regions = ["us-east-1"]
	})
}
`, tfControlPolicyResource, ControlPolicyResourceName, ControlPolicyName, ControlPolicyDescription, ControlPolicyType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyResourceName(ControlPolicyResourceName), "name", ControlPolicyName),
					resource.TestCheckResourceAttr(controlPolicyResourceName(ControlPolicyResourceName), "description", ControlPolicyDescription),
					resource.TestCheckResourceAttr(controlPolicyResourceName(ControlPolicyResourceName), "type", ControlPolicyType),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(controlPolicyResourceName(ControlPolicyResourceName), "id"),
					resource.TestCheckResourceAttrSet(controlPolicyResourceName(ControlPolicyResourceName), "parameters"),
				),
			},
			// Update and Read testing
			{
				ConfigVariables: config.Variables{
					"parameters": config.StringVariable(ControlPolicyParametersAfterUpdate),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "parameters" {
	type = string
}

resource "%s" "%s" {
	name = "%s"
	description = "%s"
	type = "%s"
	parameters = var.parameters
}
`, tfControlPolicyResource, ControlPolicyResourceName, ControlPolicyNameAfterUpdate, ControlPolicyDescription, ControlPolicyType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyResourceName(ControlPolicyResourceName), "name", ControlPolicyNameAfterUpdate),
					resource.TestCheckResourceAttr(controlPolicyResourceName(ControlPolicyResourceName), "description", ControlPolicyDescription),
					resource.TestCheckResourceAttr(controlPolicyResourceName(ControlPolicyResourceName), "type", ControlPolicyType),
					resource.TestCheckResourceAttr(controlPolicyResourceName(ControlPolicyResourceName), "parameters", ControlPolicyParametersAfterUpdate),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(controlPolicyResourceName(ControlPolicyResourceName), "id"),
				),
			},
			{
				ResourceName:      controlPolicyResourceName(ControlPolicyResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ConfigVariables: config.Variables{
					"parameters": config.StringVariable(ControlPolicyParametersAfterUpdate),
				},
			},
		},
	})
}

func controlPolicyResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfControlPolicyResource, s)
}
