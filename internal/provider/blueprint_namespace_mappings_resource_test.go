package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmBlueprintNamespaceMappings = "cm_blueprint_namespace_mappings"

	blueprintNamespaceMappingsResourceName = "blueprint_namespaces"
)

func testAccBlueprintNamespaceMappingsResourceSetup(providerId string, repoName string) string {
	return fmt.Sprintf(`
resource "cm_blueprint" "test_blueprint" {
    name = "Variable Test Blueprint"
    description = "Blueprint for testing variables"

    blueprint_vcs_info = {
        provider_id = "%s"
        repo_name = "%s"
        path = "cm/blueprint"
    }

    stack_configuration = {
        name_pattern = "{stack_name}"
        iac_type = "terraform"

        vcs_info_with_patterns = {
            provider_id = "%s"
            repo_name = "%s"
            path_pattern = "{stack_path}"
        }
    }

    substitute_parameters = [
		{
			key = "stack_name"
			description = "any name you want"
		},
		{
			key = "stack_path"
			description = "path"
		}
	]
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
`, providerId, repoName, providerId, repoName)
}

func TestAccBlueprintNamespaceMappingsResource(t *testing.T) {
	providerId := test_config.GetProviderId()
	repoName := test_config.GetRepoName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccBlueprintNamespaceMappingsResourceSetup(providerId, repoName) + fmt.Sprintf(`
resource "%s" "%s" {
  blueprint_id = cm_blueprint.test_blueprint.id
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
`, cmBlueprintNamespaceMappings, blueprintNamespaceMappingsResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "id"),
					resource.TestCheckResourceAttrSet(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "blueprint_id"),
					resource.TestCheckResourceAttr(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "namespaces.#", "3"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccBlueprintNamespaceMappingsResourceSetup(providerId, repoName) + fmt.Sprintf(`
resource "%s" "%s" {
 blueprint_id = cm_blueprint.test_blueprint.id
 namespaces = [
  {
   namespace_id = cm_namespace.test_mapping_namespace.id
  }
 ]
}
`, cmBlueprintNamespaceMappings, blueprintNamespaceMappingsResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "id"),
					resource.TestCheckResourceAttrSet(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "blueprint_id"),
					resource.TestCheckResourceAttr(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "namespaces.#", "1"),
					resource.TestCheckResourceAttrSet(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "namespaces.0.namespace_id"),
				),
			},
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmBlueprintNamespaceMappings, blueprintNamespaceMappingsResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "id"),
					resource.TestCheckResourceAttrPair(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "id", "cm_blueprint.test_blueprint", "id"),
					resource.TestCheckResourceAttrPair(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "blueprint_id", "cm_blueprint.test_blueprint", "id"),
					resource.TestCheckResourceAttrSet(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "namespaces"),
				),
			},
		},
	})
}

func blueprintNamespaceMappingsResource(s string) string {
	return fmt.Sprintf("%s.%s", cmBlueprintNamespaceMappings, s)
}
