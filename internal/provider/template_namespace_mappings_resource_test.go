package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmTemplateNamespaceMappings = "cm_template_namespace_mappings"

	templateNamespaceMappingsResourceName = "template_namespaces"
)

func testAccTemplateNamespaceMappingsResourceSetup(providerId string, repoName string) string {
	return fmt.Sprintf(`
resource "cm_template" "test_template" {
  name                = "TestTemplate"
  iac_type = "terraform"

  vcs_info = {
    provider_id = "%s"
    repo_name   = "%s"
    path        = "terraform"
  }

  policy = {
    ttl_config = {
      max_ttl = {
        type  = "days"
        value = "10"
      }
      default_ttl = {
        type  = "days"
        value = "5"
      }
    }
  }
}

resource "cm_namespace" "test_mapping_namespace" {
  name = "TestMappingNamespace"
}

resource "cm_namespace" "namespace"{
  name = "Namespace Resource"
}

resource "cm_namespace" "namespace2" {
  name = "Namespace Resource2"
}
`, providerId, repoName)
}

func TestAccTemplateNamespaceMappingsResource(t *testing.T) {
	// Test environment variables used by this function
	providerId := test_config.GetProviderId()
	repoName := test_config.GetRepoName()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccTemplateNamespaceMappingsResourceSetup(providerId, repoName) + fmt.Sprintf(`
resource "%s" "%s" {
  template_id = cm_template.test_template.id
  namespaces = [
    {
      namespace_id = cm_namespace.test_mapping_namespace.id
    },
    {
      namespace_id = cm_namespace.namespace.id
    },
    {
      namespace_id = cm_namespace.namespace2.id
    },
  ]
}
`, cmTemplateNamespaceMappings, templateNamespaceMappingsResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "id"),
					resource.TestCheckResourceAttrSet(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "template_id"),
					resource.TestCheckResourceAttr(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "namespaces.#", "3"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccTemplateNamespaceMappingsResourceSetup(providerId, repoName) + fmt.Sprintf(`
resource "%s" "%s" {
 template_id = cm_template.test_template.id
 namespaces = [
  {
   namespace_id = cm_namespace.test_mapping_namespace.id
  }
 ]
}
`, cmTemplateNamespaceMappings, templateNamespaceMappingsResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "id"),
					resource.TestCheckResourceAttrSet(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "template_id"),
					resource.TestCheckResourceAttr(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "namespaces.#", "1"),
					resource.TestCheckResourceAttrSet(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "namespaces.0.namespace_id"),
				),
			},
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmTemplateNamespaceMappings, templateNamespaceMappingsResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "id"),
					resource.TestCheckResourceAttrPair(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "id", "cm_template.test_template", "id"),
					resource.TestCheckResourceAttrPair(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "template_id", "cm_template.test_template", "id"),
					resource.TestCheckResourceAttrSet(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "namespaces"),
				),
			},
		},
	})
}

func templateNamespaceMappingsResource(s string) string {
	return fmt.Sprintf("%s.%s", cmTemplateNamespaceMappings, s)
}
