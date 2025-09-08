package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmStack = "cm_stack"

	s1ResourceName              = "stack1"
	s1IacType                   = "terraform"
	s1Name                      = "stack1"
	s1Description               = "hi"
	s1DeployOnPush              = "false"
	s1WaitForApproval           = "true"
	s1TerraformVersion          = "1.4.5"
	s1RunTriggerPatternsElement = "hi"
	s1PolicyTtlType             = "hours"
	s1PolicyTtlValue            = "1"

	s1NameAfterUpdate             = "stack2"
	s1IacTypeAfterUpdate          = "terragrunt"
	s1TerrgruntVersionAfterUpdate = "0.45.3"
)

var (
	s1ProviderId = os.Getenv("CM_TEST_PROVIDER_ID")
	s1RepoName   = os.Getenv("CM_TEST_REPO_NAME")
)

// should return 400
func TestAccStackResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccStackResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
 iac_type = "%s"
 namespace_id = cm_namespace.test_namespace.id
 name = "%s"
 description = "%s"
 deployment_behavior = {
   deploy_on_push = %s
   wait_for_approval = %s
 }
 vcs_info = {
   provider_id = "%s"
   repo_name = "%s"
 }
 iac_config = {
	terraform_version = "%s"
 }
 run_trigger = {
	patterns = ["%s"]
 }
 policy = {
	ttl_config = {
	  ttl = {
	    type = "%s"
	    value = %s
	  }
	}
 }
}
`, cmStack, s1ResourceName, s1IacType, s1Name, s1Description, s1DeployOnPush, s1WaitForApproval,
					s1ProviderId, s1RepoName, s1TerraformVersion, s1RunTriggerPatternsElement,
					s1PolicyTtlType, s1PolicyTtlValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "iac_type", s1IacType),
					resource.TestCheckResourceAttrSet(stackResourceName(s1ResourceName), "namespace_id"),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "name", s1Name),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "description", s1Description),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "deployment_behavior.deploy_on_push", s1DeployOnPush),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "deployment_behavior.wait_for_approval", s1WaitForApproval),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "vcs_info.provider_id", s1ProviderId),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "vcs_info.repo_name", s1RepoName),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "iac_config.terraform_version", s1TerraformVersion),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "run_trigger.patterns.0", s1RunTriggerPatternsElement),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "policy.ttl_config.ttl.type", s1PolicyTtlType),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "policy.ttl_config.ttl.value", s1PolicyTtlValue),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(stackResourceName(s1ResourceName), "id"),
				),
			},
			// Update and Read testing
			{
				ConfigVariables: config.Variables{
					"trigger_patterns": config.ListVariable(config.StringVariable("a"), config.StringVariable("b")),
				},
				Config: testAccStackResourceSetup() + fmt.Sprintf(`

variable "trigger_patterns" {
  type = list(string)
  default = []
}

resource "%s" "%s" {
  iac_type = "%s"
  namespace_id = cm_namespace.test_namespace.id
  name = "%s"
  deployment_behavior = {
    deploy_on_push = %s
    wait_for_approval = %s
  }
  vcs_info = {
    provider_id = "%s"
    repo_name = "%s"
  }
  iac_config = {
 	terragrunt_version = "%s"
  }
  run_trigger = length(var.trigger_patterns) == 0 ? null : {
      patterns = var.trigger_patterns
  }
   policy = {
     ttl_config = {
 	  ttl = {
 	    type = "%s"
 	    value = %s
 	  }
 	}
   }
 }
`, cmStack, s1ResourceName, s1IacTypeAfterUpdate, s1NameAfterUpdate, s1DeployOnPush, s1WaitForApproval,
					s1ProviderId, s1RepoName, s1TerrgruntVersionAfterUpdate, s1PolicyTtlType, s1PolicyTtlValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(stackResourceName(s1ResourceName), "id"),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "iac_type", s1IacTypeAfterUpdate),
					resource.TestCheckResourceAttrSet(stackResourceName(s1ResourceName), "namespace_id"),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "name", s1NameAfterUpdate),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "deployment_behavior.deploy_on_push", s1DeployOnPush),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "deployment_behavior.wait_for_approval", s1WaitForApproval),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "vcs_info.provider_id", s1ProviderId),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "vcs_info.repo_name", s1RepoName),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "iac_config.terragrunt_version", s1TerrgruntVersionAfterUpdate),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "policy.ttl_config.ttl.type", s1PolicyTtlType),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "policy.ttl_config.ttl.value", s1PolicyTtlValue),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "run_trigger.patterns.0", "a"),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "run_trigger.patterns.1", "b"),

					resource.TestCheckNoResourceAttr(stackResourceName(s1ResourceName), "description"),
					resource.TestCheckNoResourceAttr(stackResourceName(s1ResourceName), "iac_config.terraform_version"),
				),
			},
			// Update and Read testing
			{
				Config: testAccStackResourceSetup() + fmt.Sprintf(`

variable "trigger_patterns" {
  type = list(string)
  default = []
}

resource "%s" "%s" {
  iac_type = "%s"
  namespace_id = cm_namespace.test_namespace.id
  name = "%s"
  deployment_behavior = {
    deploy_on_push = %s
    wait_for_approval = %s
  }
  vcs_info = {
    provider_id = "%s"
    repo_name = "%s"
  }
  iac_config = {
 	terragrunt_version = "%s"
  }
  run_trigger = length(var.trigger_patterns) == 0 ? null : {
      patterns = var.trigger_patterns
  }
   policy = {
     ttl_config = {
 	  ttl = {
 	    type = "%s"
 	    value = %s
 	  }
 	}
   }
 }
`, cmStack, s1ResourceName, s1IacTypeAfterUpdate, s1NameAfterUpdate, s1DeployOnPush, s1WaitForApproval,
					s1ProviderId, s1RepoName, s1TerrgruntVersionAfterUpdate, s1PolicyTtlType, s1PolicyTtlValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(stackResourceName(s1ResourceName), "id"),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "iac_type", s1IacTypeAfterUpdate),
					resource.TestCheckResourceAttrSet(stackResourceName(s1ResourceName), "namespace_id"),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "name", s1NameAfterUpdate),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "deployment_behavior.deploy_on_push", s1DeployOnPush),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "deployment_behavior.wait_for_approval", s1WaitForApproval),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "vcs_info.provider_id", s1ProviderId),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "vcs_info.repo_name", s1RepoName),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "iac_config.terragrunt_version", s1TerrgruntVersionAfterUpdate),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "policy.ttl_config.ttl.type", s1PolicyTtlType),
					resource.TestCheckResourceAttr(stackResourceName(s1ResourceName), "policy.ttl_config.ttl.value", s1PolicyTtlValue),

					resource.TestCheckNoResourceAttr(stackResourceName(s1ResourceName), "run_trigger"),
					resource.TestCheckNoResourceAttr(stackResourceName(s1ResourceName), "description"),
					resource.TestCheckNoResourceAttr(stackResourceName(s1ResourceName), "iac_config.terraform_version"),
				),
			},
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmStack, s1ResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func stackResourceName(s string) string {
	return fmt.Sprintf("%s.%s", cmStack, s)
}

func testAccStackResourceSetup() string {
	return providerConfig + fmt.Sprintf(`
resource "cm_namespace" "test_namespace" {
  name = "Stack Namespace"
}
`)
}
