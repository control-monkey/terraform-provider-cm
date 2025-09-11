package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfCmControlPolicyMapping = "cm_control_policy_mappings"

	mappingResourceName = "mapping"
)

func testAccControlPolicyMappingResourceSetup() string {
	// Test environment variables used by this function
	providerId := test_config.GetProviderId()
	repoName := test_config.GetRepoName()

	return fmt.Sprintf(`
resource "cm_control_policy" "test_control_policy" {
  name = "Control Policy Unique"
  type = "aws_denied_regions"
  parameters = jsonencode({
    regions = ["us-east-1"]
  })
}

resource "cm_namespace" "dev_namespace" {
  name = "Dev"
}

resource "cm_namespace" "dev_namespace2" {
  name = "Dev2"
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

func TestAccControlPolicyMappingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccControlPolicyMappingResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  control_policy_id = cm_control_policy.test_control_policy.id
  targets = [
	{
  	  target_id         = cm_namespace.dev_namespace2.id
  	  target_type       = "namespace"
  	  enforcement_level = "%s"
	},
	{
  	  target_id         = cm_namespace.dev_namespace.id
  	  target_type       = "namespace"
  	  enforcement_level = "softMandatory"
	},
	{
  	  target_id         = cm_stack.target.id
  	  target_type       = "stack"
  	  enforcement_level = "hardMandatory"
	},
  ]
}
`, tfCmControlPolicyMapping, mappingResourceName, enforcementLevel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(controlPolicyMappingResourceName(mappingResourceName), "control_policy_id"),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "targets.#", "3"),

					resource.TestCheckResourceAttrSet(controlPolicyMappingResourceName(mappingResourceName), "id"),
					resource.TestCheckResourceAttrPair(controlPolicyMappingResourceName(mappingResourceName), "id", "cm_control_policy.test_control_policy", "id"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccControlPolicyMappingResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  control_policy_id = cm_control_policy.test_control_policy.id
  targets = [
	{
  	  target_id         = cm_namespace.dev_namespace2.id
  	  target_type       = "namespace"
  	  enforcement_level = "%s"
	},
	{
  	  target_id         = cm_namespace.dev_namespace.id
  	  target_type       = "namespace"
  	  enforcement_level = "softMandatory"
	},
  ]
}`, tfCmControlPolicyMapping, mappingResourceName, enforcementLevelAfterUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(controlPolicyMappingResourceName(mappingResourceName), "control_policy_id"),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "targets.#", "2"),

					resource.TestCheckResourceAttrSet(controlPolicyMappingResourceName(mappingResourceName), "id"),
					resource.TestCheckResourceAttrPair(controlPolicyMappingResourceName(mappingResourceName), "id", "cm_control_policy.test_control_policy", "id"),
				),
			},
			{
				ResourceName:      controlPolicyMappingResourceName(mappingResourceName),
				ImportStateVerify: true,
				ImportState:       true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(controlPolicyMappingResourceName(mappingResourceName), "control_policy_id"),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "targets.#", "2"),

					resource.TestCheckResourceAttrSet(controlPolicyMappingResourceName(mappingResourceName), "id"),
				),
			},
		},
	})
}

func controlPolicyMappingResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfCmControlPolicyMapping, s)
}
