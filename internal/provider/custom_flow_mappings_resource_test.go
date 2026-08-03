package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfCmCustomFlowMapping = "cm_custom_flow_mappings"

	cfmResourceName = "custom_flow_mapping"
)

// Two custom flows on purpose: the list endpoint returns every mapping in the organization, so the
// second flow proves this resource keeps only the mappings it owns.
func testAccCustomFlowMappingResourceSetup() string {
	// Test environment variables used by this function
	providerId := test_config.GetProviderId()
	repoName := test_config.GetRepoName()

	return fmt.Sprintf(`
resource "cm_custom_flow" "cfm_flow" {
  name       = "Custom Flow For Mapping"
  is_enabled = true
  flow_yaml  = "version: 1"
}

resource "cm_custom_flow" "cfm_other_flow" {
  name       = "Other Custom Flow"
  is_enabled = true
  flow_yaml  = "version: 1"
}

resource "cm_namespace" "cfm_namespace" {
  name = "Custom Flow Mapping Namespace"
}

resource "cm_stack" "cfm_stack" {
  iac_type     = "terraform"
  namespace_id = cm_namespace.cfm_namespace.id
  name         = "Custom Flow Mapping Stack"
  deployment_behavior = {
    deploy_on_push = false
  }
  vcs_info = {
    provider_id = "%s"
    repo_name   = "%s"
  }
}

resource "cm_custom_flow_mappings" "cfm_other_mapping" {
  custom_flow_id = cm_custom_flow.cfm_other_flow.id
  targets = [
	{
  	  target_id   = cm_namespace.cfm_namespace.id
  	  target_type = "namespace"
	},
  ]
}
`, providerId, repoName)
}

func TestAccCustomFlowMappingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create with 2 targets while another flow also has a mapping
			{
				Config: providerConfig + testAccCustomFlowMappingResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  custom_flow_id = cm_custom_flow.cfm_flow.id
  targets = [
	{
  	  target_id   = cm_namespace.cfm_namespace.id
  	  target_type = "namespace"
	},
	{
  	  target_id   = cm_stack.cfm_stack.id
  	  target_type = "stack"
	},
  ]
}
`, tfCmCustomFlowMapping, cfmResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(customFlowMappingResourceName(cfmResourceName), "custom_flow_id"),
					// 2, not 3 - the other flow's mapping must not leak into this resource
					resource.TestCheckResourceAttr(customFlowMappingResourceName(cfmResourceName), "targets.#", "2"),

					resource.TestCheckResourceAttrSet(customFlowMappingResourceName(cfmResourceName), "id"),
					resource.TestCheckResourceAttrPair(customFlowMappingResourceName(cfmResourceName), "id", "cm_custom_flow.cfm_flow", "id"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			// Replace the stack target with the organization-wide ALL namespace target
			{
				Config: providerConfig + testAccCustomFlowMappingResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  custom_flow_id = cm_custom_flow.cfm_flow.id
  targets = [
	{
  	  target_id   = cm_namespace.cfm_namespace.id
  	  target_type = "namespace"
	},
	{
  	  target_id   = "ALL"
  	  target_type = "namespace"
	},
  ]
}
`, tfCmCustomFlowMapping, cfmResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(customFlowMappingResourceName(cfmResourceName), "targets.#", "2"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			{
				ResourceName:      customFlowMappingResourceName(cfmResourceName),
				ImportStateVerify: true,
				ImportState:       true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(customFlowMappingResourceName(cfmResourceName), "custom_flow_id"),
					resource.TestCheckResourceAttr(customFlowMappingResourceName(cfmResourceName), "targets.#", "2"),

					resource.TestCheckResourceAttrSet(customFlowMappingResourceName(cfmResourceName), "id"),
				),
			},
		},
	})
}

func customFlowMappingResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfCmCustomFlowMapping, s)
}
