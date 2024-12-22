package provider

import (
	"fmt"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfDisasterRecoveryConfigurationResourceResource = "cm_disaster_recovery_configuration"
	disasterRecoveryConfigurationTfResourceName     = "disasterRecoveryConfiguration"

	disasterRecoveryConfigurationScope                    = "aws"
	disasterRecoveryConfigurationMode                     = "default"
	disasterRecoveryConfigurationIncludeManaged           = "true"
	disasterRecoveryConfigurationProviderId               = s1ProviderId
	disasterRecoveryConfigurationRepoName                 = s1RepoName
	disasterRecoveryConfigurationRepoBranch               = "main"
	disasterRecoveryConfigurationBackupStrategyGroupsJson = "[{\"vcsInfo\":{\"path\":\"a/b/c\"},\"awsQuery\":{\"region\":\"us-east-1\",\"services\":[\"AWS::EC2\"],\"resourceTypes\":[\"AWS::EC2::Instance\"],\"tags\":[{\"key\":\"Owner\",\"value\":\"Me\"}],\"excludeTags\":[{\"key\":\"Owner2\",\"value\":\"Me2\"}]}}]\n"

	disasterRecoveryConfigurationModeAfterUpdate     = "manual"
	disasterRecoveryConfigurationRepoNameAfterUpdate = "terraform"
)

// normalize json string
var disasterRecoveryConfigurationBackupStrategyGroupsJsonAfterUpdateString = "[{\"vcsInfo\":{\"path\":\"a/b/c\"},\"awsQuery\":{\"region\":\"us-east-1\",\"services\":[\"AWS::S3\"],\"resourceTypes\":[\"AWS::S3::Bucket\"],\"tags\":[{\"key\":\"Owner\",\"value\":\"Me\"}],\"excludeTags\":[{\"key\":\"Owner3\",\"value\":\"Me3\"}]}}]\n"
var disasterRecoveryConfigurationBackupStrategyGroupsJsonAfterUpdate = helpers.NormalizeJsonArrayString(disasterRecoveryConfigurationBackupStrategyGroupsJsonAfterUpdateString)

func TestAccDisasterRecoveryConfigurationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(os.Getenv("CLOUD_ACCOUNT_ID")),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "cloud_account_id" {
	type = string
}

resource "%s" "%s" {
    scope = "%s"
    cloud_account_id = var.cloud_account_id

	backup_strategy = {
    	include_managed_resources = %s
    	mode = "%s"

    	vcs_info = {
			provider_id = "%s"
			repo_name   = "%s"
			branch      = "%s"
    	}
	}
}
`, tfDisasterRecoveryConfigurationResourceResource, disasterRecoveryConfigurationTfResourceName,
					disasterRecoveryConfigurationScope, disasterRecoveryConfigurationIncludeManaged,
					disasterRecoveryConfigurationMode, disasterRecoveryConfigurationProviderId,
					disasterRecoveryConfigurationRepoName,
					disasterRecoveryConfigurationRepoBranch),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "scope", disasterRecoveryConfigurationScope),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", os.Getenv("CLOUD_ACCOUNT_ID")),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationMode),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", disasterRecoveryConfigurationProviderId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", disasterRecoveryConfigurationRepoName),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckNoResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(os.Getenv("CLOUD_ACCOUNT_ID")),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "cloud_account_id" {
	type = string
}

resource "%s" "%s" {
    scope = "%s"
    cloud_account_id = var.cloud_account_id

	backup_strategy = {
    	include_managed_resources = %s
    	mode = "%s"

    	vcs_info = {
			provider_id = "%s"
			repo_name   = "%s"
			branch      = "%s"
    	}

      groups = jsonencode([
          {
            vcsInfo = {
              path = "a/b/c"
            }

            awsQuery = {
              region = "us-east-1"
              services = ["AWS::EC2"]
              resourceTypes = ["AWS::EC2::Instance"]

              tags = [{
                key   = "Owner"
                value = "Me"
              }]
              excludeTags = [{
                key   = "Owner2"
                value = "Me2"
              }]
            }
          },
        ])
  	}
}
`, tfDisasterRecoveryConfigurationResourceResource, disasterRecoveryConfigurationTfResourceName,
					disasterRecoveryConfigurationScope, disasterRecoveryConfigurationIncludeManaged,
					disasterRecoveryConfigurationModeAfterUpdate, disasterRecoveryConfigurationProviderId,
					disasterRecoveryConfigurationRepoName,
					disasterRecoveryConfigurationRepoBranch),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "scope", disasterRecoveryConfigurationScope),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", os.Getenv("CLOUD_ACCOUNT_ID")),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationModeAfterUpdate),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", disasterRecoveryConfigurationProviderId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", disasterRecoveryConfigurationRepoName),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(os.Getenv("CLOUD_ACCOUNT_ID")),
					"groups_json":      config.StringVariable(disasterRecoveryConfigurationBackupStrategyGroupsJson),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "cloud_account_id" {
	type = string
}

variable "groups_json" {
	type = string
}

resource "%s" "%s" {
    scope = "%s"
    cloud_account_id = var.cloud_account_id

	backup_strategy = {
    	include_managed_resources = %s
    	mode = "%s"

    	vcs_info = {
			provider_id = "%s"
			repo_name   = "%s"
			branch      = "%s"
    	}

      groups = var.groups_json
  	}
}
`, tfDisasterRecoveryConfigurationResourceResource, disasterRecoveryConfigurationTfResourceName,
					disasterRecoveryConfigurationScope, disasterRecoveryConfigurationIncludeManaged,
					disasterRecoveryConfigurationModeAfterUpdate, disasterRecoveryConfigurationProviderId,
					disasterRecoveryConfigurationRepoName,
					disasterRecoveryConfigurationRepoBranch),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "scope", disasterRecoveryConfigurationScope),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", os.Getenv("CLOUD_ACCOUNT_ID")),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationModeAfterUpdate),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", disasterRecoveryConfigurationProviderId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", disasterRecoveryConfigurationRepoName),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups", disasterRecoveryConfigurationBackupStrategyGroupsJson),
				),
			},
			test_helpers.GetValidateNoDriftStep(),
			{
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(os.Getenv("CLOUD_ACCOUNT_ID")),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "cloud_account_id" {
	type = string
}

resource "%s" "%s" {
    scope = "%s"
    cloud_account_id = var.cloud_account_id

	backup_strategy = {
    	include_managed_resources = %s
    	mode = "%s"

    	vcs_info = {
			provider_id = "%s"
			repo_name   = "%s"
			branch      = "%s"
    	}

      groups = jsonencode([
		{
		  "vcsInfo": {
			"path": "a/b/c"
		  },
		  "awsQuery": {
			"region": "us-east-1",
			"services": [
			  "AWS::EC2"
			],
			"resourceTypes": [
			  "AWS::EC2::Instance"
			],
			"tags": [
			  {
				"key": "Owner",
				"value": "Me"
			  }
			],
			"excludeTags": [
			  {
				"key": "Owner2",
				"value": "Me2"
			  }
			]
		  }
		}
	  ])
  	}
}
`, tfDisasterRecoveryConfigurationResourceResource, disasterRecoveryConfigurationTfResourceName,
					disasterRecoveryConfigurationScope, disasterRecoveryConfigurationIncludeManaged,
					disasterRecoveryConfigurationModeAfterUpdate, disasterRecoveryConfigurationProviderId,
					disasterRecoveryConfigurationRepoName,
					disasterRecoveryConfigurationRepoBranch),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "scope", disasterRecoveryConfigurationScope),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", os.Getenv("CLOUD_ACCOUNT_ID")),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationModeAfterUpdate),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", disasterRecoveryConfigurationProviderId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", disasterRecoveryConfigurationRepoName),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(os.Getenv("CLOUD_ACCOUNT_ID")),
					"groups_json":      config.StringVariable(disasterRecoveryConfigurationBackupStrategyGroupsJsonAfterUpdate),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "cloud_account_id" {
	type = string
}

variable "groups_json" {
	type = string
}

resource "%s" "%s" {
    scope = "%s"
    cloud_account_id = var.cloud_account_id

	backup_strategy = {
    	include_managed_resources = %s
    	mode = "%s"

    	vcs_info = {
			provider_id = "%s"
			repo_name   = "%s"
			branch      = "%s"
    	}

      	groups = var.groups_json
  	}
}
`, tfDisasterRecoveryConfigurationResourceResource, disasterRecoveryConfigurationTfResourceName,
					disasterRecoveryConfigurationScope, disasterRecoveryConfigurationIncludeManaged,
					disasterRecoveryConfigurationModeAfterUpdate, disasterRecoveryConfigurationProviderId,
					disasterRecoveryConfigurationRepoName,
					disasterRecoveryConfigurationRepoBranch),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "scope", disasterRecoveryConfigurationScope),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", os.Getenv("CLOUD_ACCOUNT_ID")),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationModeAfterUpdate),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", disasterRecoveryConfigurationProviderId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", disasterRecoveryConfigurationRepoName),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups", disasterRecoveryConfigurationBackupStrategyGroupsJsonAfterUpdate),
				),
			},
			test_helpers.GetValidateNoDriftStep(),
			{
				ResourceName:      fmt.Sprintf("%s.%s", tfDisasterRecoveryConfigurationResourceResource, disasterRecoveryConfigurationTfResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(os.Getenv("CLOUD_ACCOUNT_ID")),
					"groups_json":      config.StringVariable(disasterRecoveryConfigurationBackupStrategyGroupsJsonAfterUpdate),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "cloud_account_id" {
	type = string
}

variable "groups_json" {
	type = string
}

resource "%s" "%s" {}
`, tfDisasterRecoveryConfigurationResourceResource, disasterRecoveryConfigurationTfResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "scope", disasterRecoveryConfigurationScope),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", os.Getenv("CLOUD_ACCOUNT_ID")),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationModeAfterUpdate),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", disasterRecoveryConfigurationProviderId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", disasterRecoveryConfigurationRepoName),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups", disasterRecoveryConfigurationBackupStrategyGroupsJsonAfterUpdate),
				),
			},
			{
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(os.Getenv("CLOUD_ACCOUNT_ID")),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "cloud_account_id" {
	type = string
}

resource "%s" "%s" {
    scope = "%s"
    cloud_account_id = var.cloud_account_id

	backup_strategy = {
    	include_managed_resources = %s
    	mode = "%s"

    	vcs_info = {
			provider_id = "%s"
			repo_name   = "%s"
			branch      = "%s"
    	}
	}
}
`, tfDisasterRecoveryConfigurationResourceResource, disasterRecoveryConfigurationTfResourceName,
					disasterRecoveryConfigurationScope, disasterRecoveryConfigurationIncludeManaged,
					disasterRecoveryConfigurationMode, disasterRecoveryConfigurationProviderId,
					disasterRecoveryConfigurationRepoNameAfterUpdate,
					disasterRecoveryConfigurationRepoBranch),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "scope", disasterRecoveryConfigurationScope),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", os.Getenv("CLOUD_ACCOUNT_ID")),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationMode),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", disasterRecoveryConfigurationProviderId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", disasterRecoveryConfigurationRepoNameAfterUpdate),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckNoResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				ResourceName:      fmt.Sprintf("%s.%s", tfDisasterRecoveryConfigurationResourceResource, disasterRecoveryConfigurationTfResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(os.Getenv("CLOUD_ACCOUNT_ID")),
				},
				Config: providerConfig + fmt.Sprintf(`
variable "cloud_account_id" {
	type = string
}

resource "%s" "%s" {}
`, tfDisasterRecoveryConfigurationResourceResource, disasterRecoveryConfigurationTfResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "scope", disasterRecoveryConfigurationScope),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", os.Getenv("CLOUD_ACCOUNT_ID")),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationMode),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", disasterRecoveryConfigurationProviderId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", disasterRecoveryConfigurationRepoNameAfterUpdate),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckNoResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups"),
				),
			},
		},
	})
}

func disasterRecoveryConfigurationResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfDisasterRecoveryConfigurationResourceResource, s)
}
