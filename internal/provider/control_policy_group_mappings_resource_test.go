package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfCmControlPolicyGroupMapping = "cm_control_policy_group_mappings"

	groupMappingResourceName = "groupMapping"
	controlPolicyGroupId     = "cmpg-fkruxsuepc"
)

func TestAccControlPolicyGroupMappingResource(t *testing.T) {
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

resource "cm_namespace" "dev_namespace2" {
  name = "Dev2"
}

variable "stack_var" {
  type = string
}

resource "%s" "%s" {
  control_policy_group_id = "%s"
  targets = [
	{
  	  target_id         = "%s"
  	  target_type       = "%s"
  	  enforcement_level = "%s"
	  override_enforcements = [
	    {
		  control_policy_id = "cmp-cop9ileumk"
		  enforcement_level = "warning"
	    },
	  ]
	},
	{
  	  target_id         = cm_namespace.dev_namespace.id
  	  target_type       = "namespace"
  	  enforcement_level = "softMandatory"
	},
	{
  	  target_id         = cm_namespace.dev_namespace2.id
  	  target_type       = "namespace"
  	  enforcement_level = "warning"
	},
	{
  	  target_id         = var.stack_var
  	  target_type       = "stack"
  	  enforcement_level = "bySeverity"
	  override_enforcements = [
	    {
		  control_policy_id = "cmp-cop9ileumk"
		  enforcement_level = "softMandatory"
	    },
	  ]
	},
  ]
}
`, tfCmControlPolicyGroupMapping, groupMappingResourceName,
					controlPolicyGroupId, targetId, targetType, enforcementLevel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyGroupMappingResourceName(groupMappingResourceName), "control_policy_group_id", controlPolicyGroupId),
					resource.TestCheckResourceAttr(controlPolicyGroupMappingResourceName(groupMappingResourceName), "targets.#", "4"),

					resource.TestCheckResourceAttrSet(controlPolicyGroupMappingResourceName(groupMappingResourceName), "id"),
					resource.TestCheckResourceAttr(controlPolicyGroupMappingResourceName(groupMappingResourceName), "id", controlPolicyGroupId),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_namespace" "dev_namespace" {
  name = "Dev"
}

resource "%s" "%s" {
  control_policy_group_id = "%s"
  targets = [
	{
  	  target_id         = "%s"
  	  target_type       = "%s"
  	  enforcement_level = "%s"
	  override_enforcements = [
	    {
		  control_policy_id = "cmp-cop9ileumk"
		  enforcement_level = "warning"
	    },
	  ]
	},
  ]
}`, tfCmControlPolicyGroupMapping, groupMappingResourceName,
					controlPolicyGroupId, targetId, targetType, enforcementLevelAfterUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyGroupMappingResourceName(groupMappingResourceName), "control_policy_group_id", controlPolicyGroupId),
					resource.TestCheckResourceAttr(controlPolicyGroupMappingResourceName(groupMappingResourceName), "targets.#", "1"),

					resource.TestCheckResourceAttrSet(controlPolicyGroupMappingResourceName(groupMappingResourceName), "id"),
					resource.TestCheckResourceAttr(controlPolicyGroupMappingResourceName(groupMappingResourceName), "id", controlPolicyGroupId),
				),
			},
			{
				ResourceName:      controlPolicyGroupMappingResourceName(groupMappingResourceName),
				ImportStateVerify: true,
				ImportState:       true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyGroupMappingResourceName(groupMappingResourceName), "control_policy_group_id", controlPolicyGroupId),
					resource.TestCheckResourceAttr(controlPolicyGroupMappingResourceName(groupMappingResourceName), "targets.#", "2"),
					resource.TestCheckResourceAttr(controlPolicyGroupMappingResourceName(groupMappingResourceName), "targets.0.override_enforcements.#", "1"),

					resource.TestCheckResourceAttrSet(controlPolicyGroupMappingResourceName(groupMappingResourceName), "id"),
					resource.TestCheckResourceAttr(controlPolicyGroupMappingResourceName(groupMappingResourceName), "id", controlPolicyGroupId),
				),
			},
		},
	})
}

func controlPolicyGroupMappingResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfCmControlPolicyGroupMapping, s)
}
