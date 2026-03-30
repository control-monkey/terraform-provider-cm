package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRunTaskDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_run_task" "ds_test" {
  name       = "Run Task DS Test"
  url        = "https://my-run-task-endpoint.example.com/callback"
  is_enabled = true
}

data "cm_run_task" "ds_test" {
  name = cm_run_task.ds_test.name
}`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cm_run_task.ds_test", "id"),
					resource.TestCheckResourceAttr("data.cm_run_task.ds_test", "name", "Run Task DS Test"),
				),
			},
		},
	})
}
