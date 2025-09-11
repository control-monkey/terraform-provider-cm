package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfCmControlPolicyGroupMapping = "cm_control_policy_group_mappings"

	groupMappingResourceName = "groupMapping"
	enforcementLevel         = "hardMandatory"

	enforcementLevelAfterUpdate = "softMandatory"
)

func testAccControlPolicyGroupMappingResourceSetup() string {
	// Test environment variables used by this function
	providerId := test_config.GetProviderId()
	repoName := test_config.GetRepoName()

	return fmt.Sprintf(`
resource "cm_namespace" "dev_namespace" {
  name = "Dev"
}

resource "cm_namespace" "dev_namespace2" {
  name = "Dev2"
}

resource "cm_namespace" "test_target_namespace" {
  name = "TestTarget"
}
resource "cm_stack" "target" {
  iac_type     = "terraform"
  namespace_id = cm_namespace.dev_namespace.id
  name         = "Stack Name"
  deployment_behavior = {
    deploy_on_push    = false
  }
  vcs_info = {
    provider_id = "%s"
    repo_name   = "%s"
  }
}
`, providerId, repoName)
}

func TestAccControlPolicyGroupMappingResource(t *testing.T) {
	// Test environment variables used by this function
	controlPolicyGroupId := test_config.GetControlPolicyGroupId()
	controlPolicyId := test_config.GetControlPolicyId()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccControlPolicyGroupMappingResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  control_policy_group_id = "%s"
  targets = [
	{
  	  target_id         = cm_namespace.test_target_namespace.id
  	  target_type       = "namespace"
  	  enforcement_level = "%s"
	  override_enforcements = [
	    {
		  control_policy_id = "%s"
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
  	  target_id         = cm_stack.target.id
  	  target_type       = "stack"
  	  enforcement_level = "bySeverity"
	  override_enforcements = [
	    {
		  control_policy_id = "%s"
		  enforcement_level = "softMandatory"
	    },
	  ]
	},
  ]
}
`, tfCmControlPolicyGroupMapping, groupMappingResourceName,
					controlPolicyGroupId, enforcementLevel, controlPolicyId, controlPolicyId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyGroupMappingResourceName(groupMappingResourceName), "control_policy_group_id", controlPolicyGroupId),
					resource.TestCheckResourceAttr(controlPolicyGroupMappingResourceName(groupMappingResourceName), "targets.#", "4"),

					resource.TestCheckResourceAttrSet(controlPolicyGroupMappingResourceName(groupMappingResourceName), "id"),
					resource.TestCheckResourceAttr(controlPolicyGroupMappingResourceName(groupMappingResourceName), "id", controlPolicyGroupId),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccControlPolicyGroupMappingResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  control_policy_group_id = "%s"
  targets = [
	{
  	  target_id         = cm_namespace.dev_namespace.id
  	  target_type       = "namespace"
  	  enforcement_level = "%s"
	  override_enforcements = [
	    {
		  control_policy_id = "%s"
		  enforcement_level = "warning"
	    },
	  ]
	},
  ]
}`, tfCmControlPolicyGroupMapping, groupMappingResourceName,
					controlPolicyGroupId, enforcementLevelAfterUpdate, controlPolicyId),
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
