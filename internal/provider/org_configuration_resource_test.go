package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmOrgConfiguration = "cm_org_configuration"
	configResourceName = "org_configuration"
)

func TestAccOrgConfigurationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  runner_config = {
    mode = "selfHosted"
    groups = ["a", "bb"]
    is_overridable = true
  }

  s3_state_files_locations = [
    {
      bucket_name    = "bucket1"
      bucket_region  = "us-east-1"
      aws_account_id = "123456789"
    },
  ]

  report_configurations = [
    {
      enabled    = true
      type       = "weeklyReport"
      recipients = {
        all_admins                 = true
        email_addresses = ["example@example.com", "example3@example.com"]
        email_addresses_to_exclude = ["example2@example.com"]
      },
    },
  ]
}
`, cmOrgConfiguration, configResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(orgConfigurationResourceName(configResourceName), "id"),
					resource.TestCheckResourceAttr(orgConfigurationResourceName(configResourceName), "runner_config.mode", "selfHosted"),
					resource.TestCheckResourceAttr(orgConfigurationResourceName(configResourceName), "runner_config.groups.#", "2"),
					resource.TestCheckResourceAttrSet(orgConfigurationResourceName(configResourceName), "s3_state_files_locations.0.bucket_name"),
					resource.TestCheckResourceAttr(orgConfigurationResourceName(configResourceName), "report_configurations.0.recipients.email_addresses.#", "2"),
					resource.TestCheckResourceAttr(orgConfigurationResourceName(configResourceName), "report_configurations.0.recipients.email_addresses_to_exclude.#", "1"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  iac_config = {
    terraform_version  = "1.5.0"
    terragrunt_version = "0.39.0"
    opentofu_version   = "1.6.0"
  }

  suppressed_resources = {
    managed_by_tags = [
      {
        key   = "Owner"
        value = "ControlMonkey"
      },
    ]
  }

  report_configurations = [
    {
      enabled    = true
      type       = "weeklyReport"
      recipients = {
        all_admins                 = true
        email_addresses_to_exclude = ["example2@example.com"]
      },
    },
  ]
}
`, cmOrgConfiguration, configResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(orgConfigurationResourceName(configResourceName), "id"),
					resource.TestCheckResourceAttrSet(orgConfigurationResourceName(configResourceName), "iac_config.terraform_version"),
					resource.TestCheckResourceAttrSet(orgConfigurationResourceName(configResourceName), "iac_config.terragrunt_version"),
					resource.TestCheckResourceAttrSet(orgConfigurationResourceName(configResourceName), "iac_config.opentofu_version"),
					resource.TestCheckResourceAttrSet(orgConfigurationResourceName(configResourceName), "suppressed_resources.managed_by_tags.0.key"),

					resource.TestCheckNoResourceAttr(orgConfigurationResourceName(configResourceName), "s3_state_files_locations"),
					resource.TestCheckNoResourceAttr(orgConfigurationResourceName(configResourceName), "runner_config"),
					resource.TestCheckNoResourceAttr(orgConfigurationResourceName(configResourceName), "report_configurations.0.recipients.email_addresses"),
					resource.TestCheckResourceAttr(orgConfigurationResourceName(configResourceName), "report_configurations.0.recipients.email_addresses_to_exclude.#", "1"),
				),
			},
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmOrgConfiguration, configResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func orgConfigurationResourceName(s string) string {
	return fmt.Sprintf("%s.%s", cmOrgConfiguration, s)
}
