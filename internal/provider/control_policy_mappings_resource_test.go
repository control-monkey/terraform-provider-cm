package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfCmControlPolicyMapping = "cm_control_policy_mappings"

	mappingResourceName = "mapping"
	controlPolicyId     = "pol-868pl1qopp"
	targetId            = "ns-x82yjdyahc"
	targetType          = "namespace"
	enforcementLevel    = "hardMandatory"

	enforcementLevelAfterUpdate = "softMandatory"
)

func TestAccControlPolicyMappingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"stack_var": config.StringVariable("stk-jtnpc6pm34"),
				},
				Config: providerConfig + fmt.Sprintf(`
resource "cm_namespace" "dev_namespace" {
  name = "Dev"
}


variable "stack_var" {
  type = string
}

resource "%s" "%s" {
  control_policy_id = "%s"
  targets = [
	{
  	  target_id         = "%s"
  	  target_type       = "%s"
  	  enforcement_level = "%s"
	},
	{
  	  target_id         = cm_namespace.dev_namespace.id
  	  target_type       = "namespace"
  	  enforcement_level = "softMandatory"
	},
	{
  	  target_id         = var.stack_var
  	  target_type       = "stack"
  	  enforcement_level = "hardMandatory"
	},
  ]
}
`, tfCmControlPolicyMapping, mappingResourceName,
					controlPolicyId, targetId, targetType, enforcementLevel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "control_policy_id", controlPolicyId),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "targets.#", "3"),

					resource.TestCheckResourceAttrSet(controlPolicyMappingResourceName(mappingResourceName), "id"),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "id", getId()),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_namespace" "dev_namespace" {
  name = "Dev"
}

resource "%s" "%s" {
  control_policy_id = "%s"
  targets = [
	{
  	  target_id         = "%s"
  	  target_type       = "%s"
  	  enforcement_level = "%s"
	},
	{
  	  target_id         = cm_namespace.dev_namespace.id
  	  target_type       = "namespace"
  	  enforcement_level = "softMandatory"
	},
  ]
}`, tfCmControlPolicyMapping, mappingResourceName,
					controlPolicyId, targetId, targetType, enforcementLevelAfterUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "control_policy_id", controlPolicyId),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "targets.#", "2"),

					resource.TestCheckResourceAttrSet(controlPolicyMappingResourceName(mappingResourceName), "id"),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "id", getId()),
				),
			},
			{
				ResourceName:      controlPolicyMappingResourceName(mappingResourceName),
				ImportStateVerify: true,
				ImportState:       true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "control_policy_id", controlPolicyId),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "targets.#", "2"),

					resource.TestCheckResourceAttrSet(controlPolicyMappingResourceName(mappingResourceName), "id"),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "id", getId()),
				),
			},
		},
	})
}

func getId() string {
	return controlPolicyId
}

func controlPolicyMappingResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfCmControlPolicyMapping, s)
}
