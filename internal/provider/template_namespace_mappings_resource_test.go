package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmTemplateNamespaceMappings = "cm_template_namespace_mappings"

	templateNamespaceMappingsResourceName = "template_namespaces"
	templateId                            = "tmpl-0mc0gph0zh"
	mappingNamespaceId                    = "ns-x82yjdyahc"
)

func TestAccTemplateNamespaceMappingsResource(t *testing.T) {
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
  template_id = "%s"
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
`, cmTemplateNamespaceMappings, templateNamespaceMappingsResourceName, templateId, mappingNamespaceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "id"),
					resource.TestCheckResourceAttrSet(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "template_id"),
					resource.TestCheckResourceAttr(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "namespaces.#", "4"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
 template_id = "%s"
 namespaces = [
  {
   namespace_id = "%s"
  }
 ]
}
`, cmTemplateNamespaceMappings, templateNamespaceMappingsResourceName, templateId, mappingNamespaceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "id"),
					resource.TestCheckResourceAttr(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "template_id", templateId),
					resource.TestCheckResourceAttr(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "namespaces.#", "1"),
					resource.TestCheckResourceAttr(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "namespaces.0.namespace_id", mappingNamespaceId),
				),
			},
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmTemplateNamespaceMappings, templateNamespaceMappingsResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "id"),
					resource.TestCheckResourceAttr(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "id", templateId),
					resource.TestCheckResourceAttr(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "template_id", templateId),
					resource.TestCheckResourceAttrSet(templateNamespaceMappingsResource(templateNamespaceMappingsResourceName), "namespaces"),
				),
			},
		},
	})
}

func templateNamespaceMappingsResource(s string) string {
	return fmt.Sprintf("%s.%s", cmTemplateNamespaceMappings, s)
}
