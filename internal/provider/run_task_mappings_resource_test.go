package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfCmRunTaskMapping = "cm_run_task_mappings"

	rtmResourceName = "run_task_mapping"
)

func testAccRunTaskMappingResourceSetup() string {
	// Test environment variables used by this function
	providerId := test_config.GetProviderId()
	repoName := test_config.GetRepoName()

	return fmt.Sprintf(`
resource "cm_run_task" "test_run_task" {
  name       = "Run Task For Mapping"
  url        = "https://example.com/callback"
  is_enabled = true
}

resource "cm_namespace" "rtm_namespace" {
  name = "Run Task Mapping Namespace"
}

resource "cm_stack" "rtm_stack" {
  iac_type     = "terraform"
  namespace_id = cm_namespace.rtm_namespace.id
  name         = "Run Task Mapping Stack"
  deployment_behavior = {
    deploy_on_push = false
  }
  vcs_info = {
    provider_id = "%s"
    repo_name   = "%s"
  }
}
`, providerId, repoName)
}

func TestAccRunTaskMappingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create with 3 targets, including the organization-wide ALL namespace target
			{
				Config: providerConfig + testAccRunTaskMappingResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  run_task_id = cm_run_task.test_run_task.id
  targets = [
	{
  	  target_id         = cm_namespace.rtm_namespace.id
  	  target_type       = "namespace"
  	  enforcement_level = "warning"
  	  stage             = "postPlan"
	},
	{
  	  target_id         = cm_stack.rtm_stack.id
  	  target_type       = "stack"
  	  enforcement_level = "softMandatory"
  	  stage             = "postPlan"
	},
	{
  	  target_id         = "ALL"
  	  target_type       = "namespace"
  	  enforcement_level = "warning"
  	  stage             = "postPlan"
	},
  ]
}
`, tfCmRunTaskMapping, rtmResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(runTaskMappingResourceName(rtmResourceName), "run_task_id"),
					resource.TestCheckResourceAttr(runTaskMappingResourceName(rtmResourceName), "targets.#", "3"),

					resource.TestCheckResourceAttrSet(runTaskMappingResourceName(rtmResourceName), "id"),
					resource.TestCheckResourceAttrPair(runTaskMappingResourceName(rtmResourceName), "id", "cm_run_task.test_run_task", "id"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			// Change enforcement_level - must be an update, not delete + create
			{
				Config: providerConfig + testAccRunTaskMappingResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  run_task_id = cm_run_task.test_run_task.id
  targets = [
	{
  	  target_id         = cm_namespace.rtm_namespace.id
  	  target_type       = "namespace"
  	  enforcement_level = "hardMandatory"
  	  stage             = "postPlan"
	},
	{
  	  target_id         = cm_stack.rtm_stack.id
  	  target_type       = "stack"
  	  enforcement_level = "softMandatory"
  	  stage             = "postPlan"
	},
	{
  	  target_id         = "ALL"
  	  target_type       = "namespace"
  	  enforcement_level = "warning"
  	  stage             = "postPlan"
	},
  ]
}
`, tfCmRunTaskMapping, rtmResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(runTaskMappingResourceName(rtmResourceName), "targets.#", "3"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			// Remove a target
			{
				Config: providerConfig + testAccRunTaskMappingResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
  run_task_id = cm_run_task.test_run_task.id
  targets = [
	{
  	  target_id         = cm_namespace.rtm_namespace.id
  	  target_type       = "namespace"
  	  enforcement_level = "hardMandatory"
  	  stage             = "postPlan"
	},
  ]
}
`, tfCmRunTaskMapping, rtmResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(runTaskMappingResourceName(rtmResourceName), "targets.#", "1"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			{
				ResourceName:      runTaskMappingResourceName(rtmResourceName),
				ImportStateVerify: true,
				ImportState:       true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(runTaskMappingResourceName(rtmResourceName), "run_task_id"),
					resource.TestCheckResourceAttr(runTaskMappingResourceName(rtmResourceName), "targets.#", "1"),

					resource.TestCheckResourceAttrSet(runTaskMappingResourceName(rtmResourceName), "id"),
				),
			},
		},
	})
}

func runTaskMappingResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfCmRunTaskMapping, s)
}
