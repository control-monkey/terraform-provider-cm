package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBlueprintDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
data "cm_blueprint" "blueprint" {
  name = "Blueprint Unique"
}`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_blueprint.blueprint", "id"),
					resource.TestCheckResourceAttr("data.cm_blueprint.blueprint", "name", "Blueprint Unique"),
				),
			},
		},
	})
}
