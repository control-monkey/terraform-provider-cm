package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfCmCustomFlow = "cm_custom_flow"

	cfResourceName = "custom_flow"

	cfName            = "tf custom flow"
	cfNameAfterUpdate = "updated tf custom flow"
)

// Multi-line YAML with a block scalar, so this also checks a heredoc round-trips without drift.
const cfFlowYaml = `version: 1
customFlows:
  run:
    steps:
      terraformInit:
        before:
          - name: say hello
            cmd: |
              echo "hello"
              echo "world"
            failureBehavior: continue
`

func TestAccCustomFlowResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name       = "%s"
  is_enabled = true
  flow_yaml  = %q
}
`, tfCmCustomFlow, cfResourceName, cfName, cfFlowYaml),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(customFlowResourceName(cfResourceName), "id"),
					resource.TestCheckResourceAttr(customFlowResourceName(cfResourceName), "name", cfName),
					resource.TestCheckResourceAttr(customFlowResourceName(cfResourceName), "is_enabled", "true"),
					resource.TestCheckResourceAttr(customFlowResourceName(cfResourceName), "flow_yaml", cfFlowYaml),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			// Change only the name. flow_yaml must survive the partial update untouched.
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name       = "%s"
  is_enabled = false
  flow_yaml  = %q
}
`, tfCmCustomFlow, cfResourceName, cfNameAfterUpdate, cfFlowYaml),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(customFlowResourceName(cfResourceName), "name", cfNameAfterUpdate),
					resource.TestCheckResourceAttr(customFlowResourceName(cfResourceName), "is_enabled", "false"),
					resource.TestCheckResourceAttr(customFlowResourceName(cfResourceName), "flow_yaml", cfFlowYaml),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			{
				ResourceName:      customFlowResourceName(cfResourceName),
				ImportStateVerify: true,
				ImportState:       true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(customFlowResourceName(cfResourceName), "id"),
					resource.TestCheckResourceAttr(customFlowResourceName(cfResourceName), "flow_yaml", cfFlowYaml),
				),
			},
		},
	})
}

func customFlowResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfCmCustomFlow, s)
}
