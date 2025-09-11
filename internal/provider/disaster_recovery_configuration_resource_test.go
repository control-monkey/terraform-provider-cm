package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"
	"github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfDisasterRecoveryConfigurationResourceResource = "cm_disaster_recovery_configuration"
	disasterRecoveryConfigurationTfResourceName     = "disasterRecoveryConfiguration"

	disasterRecoveryConfigurationScope                    = "aws"
	disasterRecoveryConfigurationMode                     = "default"
	disasterRecoveryConfigurationIncludeManaged           = "true"
	disasterRecoveryConfigurationRepoBranch               = "main"
	disasterRecoveryConfigurationBackupStrategyGroupsJson = "[{\"vcsInfo\":{\"path\":\"a/b/c\"},\"awsQuery\":{\"region\":\"us-east-1\",\"services\":[\"AWS::EC2\"],\"resourceTypes\":[\"AWS::EC2::Instance\"],\"tags\":[{\"key\":\"Owner\",\"value\":\"Me\"}],\"excludeTags\":[{\"key\":\"Owner2\",\"value\":\"Me2\"}]}}]\n"

	disasterRecoveryConfigurationModeAfterUpdate = "manual"
	repoNameAfterUpdate                          = "terraform"
)

var (
	// normalize json string
	disasterRecoveryConfigurationBackupStrategyGroupsJsonAfterUpdateString = "[{\"vcsInfo\":{\"path\":\"a/b/c\"},\"awsQuery\":{\"region\":\"us-east-1\",\"services\":[\"AWS::S3\"],\"resourceTypes\":[\"AWS::S3::Bucket\"],\"tags\":[{\"key\":\"Owner\",\"value\":\"Me\"}],\"excludeTags\":[{\"key\":\"Owner3\",\"value\":\"Me3\"}]}}]\n"
	disasterRecoveryConfigurationBackupStrategyGroupsJsonAfterUpdate       = helpers.NormalizeJsonArrayString(disasterRecoveryConfigurationBackupStrategyGroupsJsonAfterUpdateString)
)

func TestAccDisasterRecoveryConfigurationResource(t *testing.T) {
	// Test environment variables used by this function
	cloudAccountId := test_config.GetCloudAccountId()
	providerId := test_config.GetProviderId()
	repoName := test_config.GetRepoName()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(cloudAccountId),
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
					disasterRecoveryConfigurationMode, providerId,
					repoName,
					disasterRecoveryConfigurationRepoBranch),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "scope", disasterRecoveryConfigurationScope),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", cloudAccountId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationMode),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", providerId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", repoName),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckNoResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups_json"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(cloudAccountId),
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

      groups_json = jsonencode([
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
					disasterRecoveryConfigurationModeAfterUpdate, providerId,
					repoName,
					disasterRecoveryConfigurationRepoBranch),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "scope", disasterRecoveryConfigurationScope),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", cloudAccountId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationModeAfterUpdate),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", providerId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", repoName),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups_json"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(cloudAccountId),
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

      groups_json = var.groups_json
  	}
}
`, tfDisasterRecoveryConfigurationResourceResource, disasterRecoveryConfigurationTfResourceName,
					disasterRecoveryConfigurationScope, disasterRecoveryConfigurationIncludeManaged,
					disasterRecoveryConfigurationModeAfterUpdate, providerId,
					repoName,
					disasterRecoveryConfigurationRepoBranch),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "scope", disasterRecoveryConfigurationScope),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", cloudAccountId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationModeAfterUpdate),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", providerId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", repoName),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups_json"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups_json", disasterRecoveryConfigurationBackupStrategyGroupsJson),
				),
			},
			test_helpers.GetValidateNoDriftStep(),
			{
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(cloudAccountId),
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

      groups_json = jsonencode([
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
					disasterRecoveryConfigurationModeAfterUpdate, providerId,
					repoName,
					disasterRecoveryConfigurationRepoBranch),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "scope", disasterRecoveryConfigurationScope),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", cloudAccountId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationModeAfterUpdate),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", providerId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", repoName),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups_json"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(cloudAccountId),
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

      	groups_json = var.groups_json
  	}
}
`, tfDisasterRecoveryConfigurationResourceResource, disasterRecoveryConfigurationTfResourceName,
					disasterRecoveryConfigurationScope, disasterRecoveryConfigurationIncludeManaged,
					disasterRecoveryConfigurationModeAfterUpdate, providerId,
					repoName,
					disasterRecoveryConfigurationRepoBranch),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "scope", disasterRecoveryConfigurationScope),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", cloudAccountId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationModeAfterUpdate),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", providerId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", repoName),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups_json"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups_json", disasterRecoveryConfigurationBackupStrategyGroupsJsonAfterUpdate),
				),
			},
			test_helpers.GetValidateNoDriftStep(),
			{
				ResourceName:      fmt.Sprintf("%s.%s", tfDisasterRecoveryConfigurationResourceResource, disasterRecoveryConfigurationTfResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(cloudAccountId),
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
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", cloudAccountId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationModeAfterUpdate),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", providerId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", repoName),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups_json"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups_json", disasterRecoveryConfigurationBackupStrategyGroupsJsonAfterUpdate),
				),
			},
			{
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(cloudAccountId),
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
					disasterRecoveryConfigurationMode, providerId,
					repoNameAfterUpdate,
					disasterRecoveryConfigurationRepoBranch),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "id"),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "scope", disasterRecoveryConfigurationScope),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", cloudAccountId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationMode),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", providerId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", repoNameAfterUpdate),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckNoResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups_json"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				ResourceName:      fmt.Sprintf("%s.%s", tfDisasterRecoveryConfigurationResourceResource, disasterRecoveryConfigurationTfResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ConfigVariables: config.Variables{
					"cloud_account_id": config.StringVariable(cloudAccountId),
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
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "cloud_account_id", cloudAccountId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.mode", disasterRecoveryConfigurationMode),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.include_managed_resources", disasterRecoveryConfigurationIncludeManaged),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.provider_id", providerId),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.repo_name", repoNameAfterUpdate),
					resource.TestCheckResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.vcs_info.branch", disasterRecoveryConfigurationRepoBranch),
					resource.TestCheckNoResourceAttr(disasterRecoveryConfigurationResourceName(disasterRecoveryConfigurationTfResourceName), "backup_strategy.groups_json"),
				),
			},
		},
	})
}

func disasterRecoveryConfigurationResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfDisasterRecoveryConfigurationResourceResource, s)
}
