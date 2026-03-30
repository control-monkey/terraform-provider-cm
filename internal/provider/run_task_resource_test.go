package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmRunTask = "cm_run_task"

	rtResourceName = "runTask"
	runTaskName    = "Test Run Task"
	runTaskUrl     = "https://my-run-task-endpoint.example.com/callback"

	runTaskNameAfterUpdate = "Test Run Task Updated"
)

func TestAccRunTaskResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name       = "%s"
  url        = "%s"
  is_enabled = true
}
`, cmRunTask, rtResourceName, runTaskName, runTaskUrl),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(runTaskResourceName(rtResourceName), "id"),
					resource.TestCheckResourceAttr(runTaskResourceName(rtResourceName), "name", runTaskName),
					resource.TestCheckResourceAttr(runTaskResourceName(rtResourceName), "url", runTaskUrl),
					resource.TestCheckResourceAttr(runTaskResourceName(rtResourceName), "is_enabled", "true"),
				),
			},
			test_helpers.GetValidateNoDriftStep(),
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name       = "%s"
  url        = "%s"
  is_enabled = false
}
`, cmRunTask, rtResourceName, runTaskNameAfterUpdate, runTaskUrl),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(runTaskResourceName(rtResourceName), "id"),
					resource.TestCheckResourceAttr(runTaskResourceName(rtResourceName), "name", runTaskNameAfterUpdate),
					resource.TestCheckResourceAttr(runTaskResourceName(rtResourceName), "is_enabled", "false"),
				),
			},
			test_helpers.GetValidateNoDriftStep(),
			{
				ResourceName:            fmt.Sprintf("%s.%s", cmRunTask, rtResourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"hmac_key"},
			},
		},
	})
}

func runTaskResourceName(s string) string {
	return fmt.Sprintf("%s.%s", cmRunTask, s)
}
