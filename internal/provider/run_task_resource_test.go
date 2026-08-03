package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const (
	cmRunTask = "cm_run_task"

	rtResourceName = "runTask"
	runTaskName    = "Test Run Task"
	runTaskUrl     = "https://my-run-task-endpoint.example.com/callback"

	runTaskNameAfterUpdate = "Test Run Task Updated"

	woRtResourceName     = "runTaskWriteOnly"
	runTaskWriteOnlyName = "Test Run Task Write Only"
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

// TestAccRunTaskWriteOnlyHmacKeyResource covers hmac_key_wo and hmac_key_wo_version.
// Skips below Terraform 1.11, which is the minimum for write-only arguments.
func TestAccRunTaskWriteOnlyHmacKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name       = "%s"
  url        = "%s"
  is_enabled = true

  hmac_key_wo         = "first-write-only-secret"
  hmac_key_wo_version = 1
}
`, cmRunTask, woRtResourceName, runTaskWriteOnlyName, runTaskUrl),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(runTaskResourceName(woRtResourceName), "id"),
					// the secret must not be persisted
					resource.TestCheckNoResourceAttr(runTaskResourceName(woRtResourceName), "hmac_key_wo"),
					resource.TestCheckNoResourceAttr(runTaskResourceName(woRtResourceName), "hmac_key"),
					// but the API must have received it
					resource.TestCheckResourceAttr(runTaskResourceName(woRtResourceName), "is_hmac_key_configured", "true"),
					resource.TestCheckResourceAttr(runTaskResourceName(woRtResourceName), "hmac_key_wo_version", "1"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),

			// Changing the key alone must produce no diff - the version is the only trigger
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name       = "%s"
  url        = "%s"
  is_enabled = true

  hmac_key_wo         = "changed-but-version-not-bumped"
  hmac_key_wo_version = 1
}
`, cmRunTask, woRtResourceName, runTaskWriteOnlyName, runTaskUrl),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},

			// Bumping the version sends the new key as an in-place update
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name       = "%s"
  url        = "%s"
  is_enabled = true

  hmac_key_wo         = "second-write-only-secret"
  hmac_key_wo_version = 2
}
`, cmRunTask, woRtResourceName, runTaskWriteOnlyName, runTaskUrl),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(runTaskResourceName(woRtResourceName), "hmac_key_wo_version", "2"),
					resource.TestCheckNoResourceAttr(runTaskResourceName(woRtResourceName), "hmac_key_wo"),
					resource.TestCheckResourceAttr(runTaskResourceName(woRtResourceName), "is_hmac_key_configured", "true"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
		},
	})
}
