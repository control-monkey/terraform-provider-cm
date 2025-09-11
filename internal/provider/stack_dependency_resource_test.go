package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmStackDependency = "cm_stack_dependency"

	sdResourceName       = "dep"
	sdStackName          = "test-stack-dependency-source"
	sdDependsOnStackName = "test-stack-dependency-target"
	sdTriggerOption      = "always"
	sdOutputName         = "db_endpoint"
	sdInputName          = "db_endpoint"
	sdIncludeSensitive   = "false"

	sdTriggerOptionUpdated = "onReferenceValueChange"
	sdOutputNameUpdated    = "vpc_id"
	sdInputNameUpdated     = "vpc_id"
)

func TestAccStackDependencyResource(t *testing.T) {
	// Test environment variables used by this function
	providerId := test_config.GetProviderId()
	repoName := test_config.GetRepoName()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				ConfigVariables: config.Variables{
					"trigger_option": config.StringVariable(sdTriggerOption),
					"output_name":    config.StringVariable(sdOutputName),
					"input_name":     config.StringVariable(sdInputName),
				},
				Config: testAccStackDependencyResourceSetup(providerId, repoName) + fmt.Sprintf(`
variable "trigger_option" {
  type = string
}

variable "output_name" {
  type = string
}

variable "input_name" {
  type = string
}

resource "%s" "%s" {
  stack_id            = cm_stack.source.id
  depends_on_stack_id = cm_stack.target.id
  trigger_option      = var.trigger_option

  references = [
    {
      output_of_stack_to_depend_on = var.output_name
      input_for_stack              = var.input_name
      include_sensitive_output     = %s
    }
  ]
}
`, cmStackDependency, sdResourceName, sdIncludeSensitive),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(stackDependencyResourceName(sdResourceName), "stack_id", "cm_stack.source", "id"),
					resource.TestCheckResourceAttrPair(stackDependencyResourceName(sdResourceName), "depends_on_stack_id", "cm_stack.target", "id"),
					resource.TestCheckResourceAttr(stackDependencyResourceName(sdResourceName), "trigger_option", sdTriggerOption),
					resource.TestCheckResourceAttr(stackDependencyResourceName(sdResourceName), "references.0.output_of_stack_to_depend_on", sdOutputName),
					resource.TestCheckResourceAttr(stackDependencyResourceName(sdResourceName), "references.0.input_for_stack", sdInputName),
					resource.TestCheckResourceAttr(stackDependencyResourceName(sdResourceName), "references.0.include_sensitive_output", sdIncludeSensitive),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(stackDependencyResourceName(sdResourceName), "id"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			// Update and Read testing
			{
				ConfigVariables: config.Variables{
					"trigger_option": config.StringVariable(sdTriggerOptionUpdated),
					"output_name":    config.StringVariable(sdOutputNameUpdated),
					"input_name":     config.StringVariable(sdInputNameUpdated),
				},
				Config: testAccStackDependencyResourceSetup(providerId, repoName) + fmt.Sprintf(`
variable "trigger_option" {
  type = string
}

variable "output_name" {
  type = string
}

variable "input_name" {
  type = string
}

resource "%s" "%s" {
  stack_id            = cm_stack.source.id
  depends_on_stack_id = cm_stack.target.id
  trigger_option      = var.trigger_option

  references = [
    {
      output_of_stack_to_depend_on = var.output_name
      input_for_stack              = var.input_name
    }
  ]
}
`, cmStackDependency, sdResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(stackDependencyResourceName(sdResourceName), "stack_id", "cm_stack.source", "id"),
					resource.TestCheckResourceAttrPair(stackDependencyResourceName(sdResourceName), "depends_on_stack_id", "cm_stack.target", "id"),
					resource.TestCheckResourceAttr(stackDependencyResourceName(sdResourceName), "trigger_option", sdTriggerOptionUpdated),
					resource.TestCheckResourceAttr(stackDependencyResourceName(sdResourceName), "references.0.output_of_stack_to_depend_on", sdOutputNameUpdated),
					resource.TestCheckResourceAttr(stackDependencyResourceName(sdResourceName), "references.0.input_for_stack", sdInputNameUpdated),
					resource.TestCheckNoResourceAttr(stackDependencyResourceName(sdResourceName), "references.0.include_sensitive_output"),
					resource.TestCheckResourceAttrSet(stackDependencyResourceName(sdResourceName), "id"),
				),
			},
			// Import testing
			{
				ConfigVariables: config.Variables{
					"trigger_option": config.StringVariable(sdTriggerOptionUpdated),
					"output_name":    config.StringVariable(sdOutputNameUpdated),
					"input_name":     config.StringVariable(sdInputNameUpdated),
				},
				ResourceName:      fmt.Sprintf("%s.%s", cmStackDependency, sdResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func stackDependencyResourceName(s string) string {
	return fmt.Sprintf("%s.%s", cmStackDependency, s)
}

func testAccStackDependencyResourceSetup(providerId string, repoName string) string {
	return providerConfig + fmt.Sprintf(`
resource "cm_namespace" "test_namespace" {
  name = "Stack Dependency Test Namespace"
}

resource "cm_stack" "source" {
  iac_type     = "terraform"
  namespace_id = cm_namespace.test_namespace.id
  name         = "%s"
  deployment_behavior = {
    deploy_on_push    = false
    wait_for_approval = false
  }
  vcs_info = {
    provider_id = "%s"
    repo_name   = "%s"
  }
  iac_config = {
    terraform_version = "1.5.0"
  }
}

resource "cm_stack" "target" {
  iac_type     = "terraform"
  namespace_id = cm_namespace.test_namespace.id
  name         = "%s"
  deployment_behavior = {
    deploy_on_push    = false
    wait_for_approval = false
  }
  vcs_info = {
    provider_id = "%s"
    repo_name   = "%s"
  }
  iac_config = {
    terraform_version = "1.5.0"
  }
}
`, sdStackName, providerId, repoName, sdDependsOnStackName, providerId, repoName)
}
