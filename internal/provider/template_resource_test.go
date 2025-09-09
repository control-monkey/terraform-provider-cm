package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmTemplate = "cm_template"

	t1ResourceName          = "template"
	t1Name                  = "Dev Self-Service Template"
	t1IacType               = "terraform"
	t1Description           = "Self service on Dev environment for developers"
	t1PolicyMaxTtlType      = "days"
	t1PolicyMaxTtlValue     = "2"
	t1PolicyDefaultTtlType  = "hours"
	t1PolicyDefaultTtlValue = "3"

	t1PolicyDefaultTtlValueAfterUpdate = "1"
	t1NameAfterUpdate                  = "Dev Self-Service Template After Update"
	t1IacTypeAfterUpdate               = "terragrunt"
)

var (
	t1ProviderId = os.Getenv("CM_TEST_PROVIDER_ID")
	t1RepoName   = os.Getenv("CM_TEST_REPO_NAME")
)

func TestAccTemplateResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
 name = "%s"
 iac_type = "%s"
 description = "%s"

 vcs_info = {
   provider_id = "%s"
   repo_name = "%s"
 }

 policy = {
	ttl_config = {
	  max_ttl = {
	    type = "%s"
	    value = %s
	  }
	  default_ttl = {
	    type = "%s"
	    value = %s
	  }
	}
 }
}
`, cmTemplate, t1ResourceName, t1Name, t1IacType, t1Description,
					t1ProviderId, t1RepoName, t1PolicyMaxTtlType, t1PolicyMaxTtlValue, t1PolicyDefaultTtlType, t1PolicyDefaultTtlValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "name", t1Name),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "iac_type", t1IacType),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "description", t1Description),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "vcs_info.provider_id", t1ProviderId),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "vcs_info.repo_name", t1RepoName),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "policy.ttl_config.max_ttl.type", t1PolicyMaxTtlType),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "policy.ttl_config.max_ttl.value", t1PolicyMaxTtlValue),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "policy.ttl_config.default_ttl.type", t1PolicyDefaultTtlType),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "policy.ttl_config.default_ttl.value", t1PolicyDefaultTtlValue),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(templateResourceName(t1ResourceName), "id"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
 name = "%s"
 iac_type = "%s"

 vcs_info = {
   provider_id = "%s"
   repo_name = "%s"
 }

 policy = {
	ttl_config = {
	  max_ttl = {
	    type = "%s"
	    value = %s
	  }
	  default_ttl = {
	    type = "%s"
	    value = %s
	  }
	}
 }
}
`, cmTemplate, t1ResourceName, t1NameAfterUpdate, t1IacTypeAfterUpdate,
					t1ProviderId, t1RepoName, t1PolicyMaxTtlType, t1PolicyMaxTtlValue, t1PolicyDefaultTtlType, t1PolicyDefaultTtlValueAfterUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(templateResourceName(t1ResourceName), "id"),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "name", t1NameAfterUpdate),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "iac_type", t1IacTypeAfterUpdate),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "vcs_info.provider_id", t1ProviderId),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "vcs_info.repo_name", t1RepoName),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "policy.ttl_config.max_ttl.type", t1PolicyMaxTtlType),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "policy.ttl_config.max_ttl.value", t1PolicyMaxTtlValue),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "policy.ttl_config.default_ttl.type", t1PolicyDefaultTtlType),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "policy.ttl_config.default_ttl.value", t1PolicyDefaultTtlValueAfterUpdate),

					resource.TestCheckNoResourceAttr(templateResourceName(t1ResourceName), "description"),
				),
			},
			// Test iac_config and runner_config
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
 name = "%s"
 iac_type = "%s"

 vcs_info = {
   provider_id = "%s"
   repo_name = "%s"
 }

 policy = {
	ttl_config = {
	  max_ttl = {
	    type = "%s"
	    value = %s
	  }
	  default_ttl = {
	    type = "%s"
	    value = %s
	  }
	}
 }

 iac_config = {
   terraform_version = "1.5.0"
   terragrunt_version = "0.45.3"
 }

 runner_config = {
   mode = "managed"
 }
}
`, cmTemplate, t1ResourceName, t1NameAfterUpdate, t1IacTypeAfterUpdate,
					t1ProviderId, t1RepoName, t1PolicyMaxTtlType, t1PolicyMaxTtlValue, t1PolicyDefaultTtlType, t1PolicyDefaultTtlValueAfterUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(templateResourceName(t1ResourceName), "id"),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "name", t1NameAfterUpdate),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "iac_type", t1IacTypeAfterUpdate),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "vcs_info.provider_id", t1ProviderId),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "vcs_info.repo_name", t1RepoName),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "policy.ttl_config.max_ttl.type", t1PolicyMaxTtlType),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "policy.ttl_config.max_ttl.value", t1PolicyMaxTtlValue),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "policy.ttl_config.default_ttl.type", t1PolicyDefaultTtlType),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "policy.ttl_config.default_ttl.value", t1PolicyDefaultTtlValueAfterUpdate),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "iac_config.terraform_version", "1.5.0"),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "iac_config.terragrunt_version", "0.45.3"),
					resource.TestCheckResourceAttr(templateResourceName(t1ResourceName), "runner_config.mode", "managed"),

					resource.TestCheckNoResourceAttr(templateResourceName(t1ResourceName), "description"),
					resource.TestCheckNoResourceAttr(templateResourceName(t1ResourceName), "iac_config.opentofu_version"),
					resource.TestCheckNoResourceAttr(templateResourceName(t1ResourceName), "runner_config.groups"),
				),
			},
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmTemplate, t1ResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func templateResourceName(s string) string {
	return fmt.Sprintf("%s.%s", cmTemplate, s)
}
