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
	cmStackDiscoveryConfiguration = "cm_stack_discovery_configuration"

	sdcResourceName     = "test"
	sdcName             = "Test Stack Discovery Configuration"
	sdcDescription      = "Test configuration for stack auto-discovery"
	sdcBranch           = "main"
	sdcIacType          = "terraform"
	sdcTerraformVersion = "1.4.5"
	sdcRunnerMode       = "managed"

	sdcNameUpdated             = "Updated Stack Discovery Configuration"
	sdcDescriptionUpdated      = "Updated test configuration"
	sdcBranchUpdated           = "develop"
	sdcTerraformVersionUpdated = "1.5.0"
)

func TestAccStackDiscoveryConfigurationResource(t *testing.T) {
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
					"provider_id":       config.StringVariable(providerId),
					"repo_name":         config.StringVariable(repoName),
					"terraform_version": config.StringVariable(sdcTerraformVersion),
				},
				Config: testAccStackDiscoveryConfigurationResourceSetup() + fmt.Sprintf(`
variable "provider_id" {
  type = string
}

variable "repo_name" {
  type = string
}

variable "terraform_version" {
  type = string
}

resource "%s" "%s" {
  name         = "%s"
  namespace_id = cm_namespace.test_namespace.id
  description  = "%s"

  vcs_patterns = [
    {
      provider_id         = var.provider_id
      repo_name           = var.repo_name
      path_patterns       = ["environments/*/**", "modules/*/**"]
      branch              = "%s"
    }
  ]

  stack_config = {
    iac_type = "%s"
    
    deployment_behavior = {
      deploy_on_push    = false
    }

    run_trigger = {
      patterns = ["**/*.tf"]
    }

    iac_config = {
      terraform_version = var.terraform_version
    }

    runner_config = {
      mode = "%s"
    }

    auto_sync = {
      deploy_when_drift_detected = true
    }
  }
}
`, cmStackDiscoveryConfiguration, sdcResourceName, sdcName, sdcDescription, sdcBranch, sdcIacType, sdcRunnerMode),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "name", sdcName),
					resource.TestCheckResourceAttrSet(stackDiscoveryConfigurationResourceName(sdcResourceName), "namespace_id"),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "description", sdcDescription),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "vcs_patterns.0.provider_id", providerId),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "vcs_patterns.0.repo_name", repoName),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "vcs_patterns.0.branch", sdcBranch),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "vcs_patterns.0.path_patterns.0", "environments/*/**"),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "vcs_patterns.0.path_patterns.1", "modules/*/**"),
					resource.TestCheckNoResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "vcs_patterns.0.exclude_path_patterns"),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "stack_config.iac_type", sdcIacType),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "stack_config.deployment_behavior.deploy_on_push", "false"),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "stack_config.iac_config.terraform_version", sdcTerraformVersion),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "stack_config.runner_config.mode", sdcRunnerMode),
					resource.TestCheckResourceAttrSet(stackDiscoveryConfigurationResourceName(sdcResourceName), "stack_config.auto_sync.deploy_when_drift_detected"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(stackDiscoveryConfigurationResourceName(sdcResourceName), "id"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			// Update and Read testing
			{
				ConfigVariables: config.Variables{
					"provider_id":       config.StringVariable(providerId),
					"repo_name":         config.StringVariable(repoName),
					"terraform_version": config.StringVariable(sdcTerraformVersionUpdated),
				},
				Config: testAccStackDiscoveryConfigurationResourceSetup() + fmt.Sprintf(`
variable "provider_id" {
  type = string
}

variable "repo_name" {
  type = string
}

variable "terraform_version" {
  type = string
}

resource "%s" "%s" {
  name         = "%s"
  namespace_id = cm_namespace.test_namespace.id
  description  = "%s"

  vcs_patterns = [
    {
      provider_id         = var.provider_id
      repo_name           = var.repo_name
      path_patterns       = ["apps/*/terraform/**", "services/*/infra/**"]
      exclude_path_patterns = ["**/node_modules/**", "**/vendor/**"]
      branch              = "%s"
    }
  ]

  stack_config = {
    iac_type = "%s"
    
    deployment_behavior = {
      deploy_on_push    = false
    }

    run_trigger = {
      patterns = ["**/*.tf", "**/*.tfvars"]
    }

    iac_config = {
      terraform_version = var.terraform_version
    }

    runner_config = {
      mode = "%s"
    }
  }
}
`, cmStackDiscoveryConfiguration, sdcResourceName, sdcNameUpdated, sdcDescriptionUpdated, sdcBranchUpdated, sdcIacType, sdcRunnerMode),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "name", sdcNameUpdated),
					resource.TestCheckResourceAttrSet(stackDiscoveryConfigurationResourceName(sdcResourceName), "namespace_id"),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "description", sdcDescriptionUpdated),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "vcs_patterns.0.branch", sdcBranchUpdated),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "vcs_patterns.0.path_patterns.0", "apps/*/terraform/**"),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "vcs_patterns.0.path_patterns.1", "services/*/infra/**"),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "vcs_patterns.0.exclude_path_patterns.#", "2"),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "stack_config.deployment_behavior.deploy_on_push", "false"),
					resource.TestCheckResourceAttr(stackDiscoveryConfigurationResourceName(sdcResourceName), "stack_config.iac_config.terraform_version", sdcTerraformVersionUpdated),
					resource.TestCheckResourceAttrSet(stackDiscoveryConfigurationResourceName(sdcResourceName), "id"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			// Import testing
			{
				ConfigVariables: config.Variables{
					"provider_id":       config.StringVariable(providerId),
					"repo_name":         config.StringVariable(repoName),
					"terraform_version": config.StringVariable(sdcTerraformVersionUpdated),
				},
				ResourceName:      fmt.Sprintf("%s.%s", cmStackDiscoveryConfiguration, sdcResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func stackDiscoveryConfigurationResourceName(s string) string {
	return fmt.Sprintf("%s.%s", cmStackDiscoveryConfiguration, s)
}

func testAccStackDiscoveryConfigurationResourceSetup() string {
	return providerConfig + fmt.Sprintf(`
resource "cm_namespace" "test_namespace" {
  name = "Stack Conviguration Namespace"
}
`)
}
