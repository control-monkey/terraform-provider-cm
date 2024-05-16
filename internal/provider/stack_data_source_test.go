package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStackDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				ConfigVariables: config.Variables{
					"stack_name": config.StringVariable("Stack Unique"),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "stack_name" {
  type = string
}

data "cm_stack" "stack" {
  name = var.stack_name
}`),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_stack.stack", "id"),
					resource.TestCheckResourceAttr("data.cm_stack.stack", "name", "Stack Unique"),
				),
			},
		},
	})
}
