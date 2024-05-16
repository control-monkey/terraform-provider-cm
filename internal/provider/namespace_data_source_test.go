package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNamespaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_namespace" "namespace"{
  name = "Namespace Unique"
}

data "cm_namespace" "namespace" {
  name = cm_namespace.namespace.name
}`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_namespace.namespace", "id"),
					resource.TestCheckResourceAttr("data.cm_namespace.namespace", "name", "Namespace Unique"),
				),
			},
		},
	})
}
