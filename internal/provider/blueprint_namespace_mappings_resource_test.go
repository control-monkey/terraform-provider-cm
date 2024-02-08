package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmBlueprintNamespaceMappings = "cm_blueprint_namespace_mappings"

	blueprintNamespaceMappingsResourceName = "blueprint_namespaces"
	blueprintId                            = "blp-1vutl5aqo2"
	mappingBlueprintNamespaceId            = "ns-x82yjdyahc"
)

func TestAccBlueprintNamespaceMappingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"namespace_var": config.StringVariable("ns-b5ol64210x"),
				},
				Config: providerConfig + fmt.Sprintf(`
resource "cm_namespace" "namespace"{
  name = "Namespace Resource"
}

resource "cm_namespace" "namespace2" {
  name = "Namespace Resource2"
}

variable "namespace_var" {
  type = string
}

resource "%s" "%s" {
  blueprint_id = "%s"
  namespaces = [
    {
      namespace_id = "%s"
    },
    {
      namespace_id = var.namespace_var
    },
    {
      namespace_id = cm_namespace.namespace.id
    },
    {
      namespace_id = cm_namespace.namespace2.id
    },
  ]
}
`, cmBlueprintNamespaceMappings, blueprintNamespaceMappingsResourceName, blueprintId, mappingBlueprintNamespaceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "id"),
					resource.TestCheckResourceAttrSet(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "blueprint_id"),
					resource.TestCheckResourceAttr(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "namespaces.#", "4"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
 blueprint_id = "%s"
 namespaces = [
  {
   namespace_id = "%s"
  }
 ]
}
`, cmBlueprintNamespaceMappings, blueprintNamespaceMappingsResourceName, blueprintId, mappingBlueprintNamespaceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "id"),
					resource.TestCheckResourceAttr(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "blueprint_id", blueprintId),
					resource.TestCheckResourceAttr(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "namespaces.#", "1"),
					resource.TestCheckResourceAttr(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "namespaces.0.namespace_id", mappingBlueprintNamespaceId),
				),
			},
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmBlueprintNamespaceMappings, blueprintNamespaceMappingsResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "id"),
					resource.TestCheckResourceAttr(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "id", blueprintId),
					resource.TestCheckResourceAttr(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "blueprint_id", blueprintId),
					resource.TestCheckResourceAttrSet(blueprintNamespaceMappingsResource(blueprintNamespaceMappingsResourceName), "namespaces"),
				),
			},
		},
	})
}

func blueprintNamespaceMappingsResource(s string) string {
	return fmt.Sprintf("%s.%s", cmBlueprintNamespaceMappings, s)
}
