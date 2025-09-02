package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"
	"github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfBlueprintResourceResource = "cm_blueprint"
	blueprintTfResourceName     = "blueprint"

	blueprintName                      = "Test Blueprint"
	blueprintDescription               = "Description"
	blueprintProviderId                = s1ProviderId
	blueprintRepoName                  = s1RepoName
	blueprintRepoPath                  = "cm/blueprint"
	blueprintStackConfigurationIacType = "terraform"

	blueprintNameAfterUpdate = "updated name"
)

func TestAccBlueprintResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
    name = "%s"
    description = "%s"

    blueprint_vcs_info = {
        provider_id = "%s"
        repo_name = "%s"
        path = "%s"
    }

    stack_configuration = {
        name_pattern = "{stack_name}"
        iac_type = "%s"

        vcs_info_with_patterns = {
            provider_id = "%s"
            repo_name = "%s"
            path_pattern = "{stack_path}"
        }
    }

    substitute_parameters = [
		{
			key = "stack_name"
			description = "any name you want"
		},
		{
			key = "stack_path"
			description = "path"
		}
	]
}
`, tfBlueprintResourceResource, blueprintTfResourceName, blueprintName, blueprintDescription, blueprintProviderId,
					blueprintRepoName, blueprintRepoPath, blueprintStackConfigurationIacType, blueprintProviderId, blueprintRepoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(blueprintResourceName(blueprintTfResourceName), "id"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "name", blueprintName),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "description", blueprintDescription),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "blueprint_vcs_info.provider_id", blueprintProviderId),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "blueprint_vcs_info.repo_name", blueprintRepoName),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "blueprint_vcs_info.path", blueprintRepoPath),
					resource.TestCheckNoResourceAttr(blueprintResourceName(blueprintTfResourceName), "blueprint_vcs_info.branch"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.name_pattern", "{stack_name}"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.iac_type", blueprintStackConfigurationIacType),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.vcs_info_with_patterns.provider_id", blueprintProviderId),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.vcs_info_with_patterns.repo_name", blueprintRepoName),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.vcs_info_with_patterns.path_pattern", "{stack_path}"),
					resource.TestCheckNoResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.vcs_info_with_patterns.branch_pattern"),
					resource.TestCheckNoResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.deployment_approval_policy"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.#", "2"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.0.key", "stack_name"),
					resource.TestCheckResourceAttrSet(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.0.description"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.1.key", "stack_path"),
					resource.TestCheckResourceAttrSet(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.1.description"),
					resource.TestCheckNoResourceAttr(blueprintResourceName(blueprintTfResourceName), "skip_plan_on_stack_initialization"),
					resource.TestCheckNoResourceAttr(blueprintResourceName(blueprintTfResourceName), "auto_approve_apply_on_initialization"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
    name = "%s"
    description = "%s"

    blueprint_vcs_info = {
        provider_id = "%s"
        repo_name = "%s"
        path = "%s"
    }

    stack_configuration = {
        name_pattern = "{stack_name}"
        iac_type = "%s"

        vcs_info_with_patterns = {
            provider_id = "%s"
            repo_name = "%s"
            path_pattern = "{stack_path}"
        }

        run_trigger = {
            patterns = ["services/{stack_name}/**", "common/**"]
            exclude_patterns = ["**/*.md"]
        }

        iac_config = {
            terraform_version   = "1.5.9"
            opentofu_version    = "1.6.0"
            var_files           = ["vars/common.tfvars"]
        }

        auto_sync = {
            deploy_when_drift_detected = true
        }
    }

    substitute_parameters = [
		{
			key = "stack_name"
			description = "any name you want"
		},
		{
			key = "stack_path"
			description = "path"
		}
	]
}
`, tfBlueprintResourceResource, blueprintTfResourceName, blueprintName, blueprintDescription, blueprintProviderId,
					blueprintRepoName, blueprintRepoPath, blueprintStackConfigurationIacType, blueprintProviderId, blueprintRepoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(blueprintResourceName(blueprintTfResourceName), "id"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.run_trigger.patterns.#", "2"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.run_trigger.exclude_patterns.#", "1"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.iac_config.terraform_version", "1.5.9"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.iac_config.opentofu_version", "1.6.0"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.iac_config.var_files.#", "1"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.auto_sync.deploy_when_drift_detected", "true"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				ConfigVariables: config.Variables{
					"lte_value": config.StringVariable("50"),
					"gte_value": config.StringVariable("5"),
				},
				Config: providerConfig + fmt.Sprintf(`
resource "cm_team" "team1" {
	name = "Developers"
}

resource "cm_team" "team2" {
	name = "QA"
}

variable "lte_value" {
	type = string
}

variable "gte_value" {
	type = string
}

resource "%s" "%s" {
    name = "%s"
    description = "%s"

    blueprint_vcs_info = {
        provider_id = "%s"
        repo_name = "%s"
        path = "%s"
        branch = "main"
    }

    stack_configuration = {
        name_pattern = "{stack_name}"
        iac_type = "%s"

        vcs_info_with_patterns = {
            provider_id = "%s"
            repo_name = "%s"
            path_pattern = "{stack_path}"
            branch_pattern = "test"
        }

		deployment_approval_policy = {
			rules = [
				{
					type = "requireTeamsApproval"
					parameters = jsonencode({
						teams = [cm_team.team1.id, cm_team.team2.id]
					})
				},
				{
					type = "requireTwoApprovals"
				},
			]
		}
    }

	substitute_parameters = [
		{
			key = "stack_name"
			description = "any name you want"
		},
		{
			key = "stack_path"
			description = "path"
		},
		{
			key = "paramNumber"
			description = "description"
			value_conditions = [
				{
					operator = "lt"
					value    = var.lte_value
				},
				{
					operator = "gt"
					value    = var.gte_value
				 },
			]
		}
	]

    skip_plan_on_stack_initialization = true
    auto_approve_apply_on_initialization = true
}
`, tfBlueprintResourceResource, blueprintTfResourceName, blueprintNameAfterUpdate, blueprintDescription, blueprintProviderId,
					blueprintRepoName, blueprintRepoPath, blueprintStackConfigurationIacType, blueprintProviderId, blueprintRepoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(blueprintResourceName(blueprintTfResourceName), "id"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "name", blueprintNameAfterUpdate),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "description", blueprintDescription),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "blueprint_vcs_info.provider_id", blueprintProviderId),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "blueprint_vcs_info.repo_name", blueprintRepoName),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "blueprint_vcs_info.path", blueprintRepoPath),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "blueprint_vcs_info.branch", "main"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.name_pattern", "{stack_name}"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.iac_type", blueprintStackConfigurationIacType),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.vcs_info_with_patterns.provider_id", blueprintProviderId),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.vcs_info_with_patterns.repo_name", blueprintRepoName),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.vcs_info_with_patterns.path_pattern", "{stack_path}"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.vcs_info_with_patterns.branch_pattern", "test"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.deployment_approval_policy.rules.#", "2"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.#", "3"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.0.key", "stack_name"),
					resource.TestCheckResourceAttrSet(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.0.description"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.1.key", "stack_path"),
					resource.TestCheckResourceAttrSet(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.1.description"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.2.key", "paramNumber"),
					resource.TestCheckResourceAttrSet(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.2.description"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.2.value_conditions.#", "2"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "skip_plan_on_stack_initialization", "true"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "auto_approve_apply_on_initialization", "true"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
    name = "%s"
    description = "%s"

    blueprint_vcs_info = {
        provider_id = "%s"
        repo_name = "%s"
        path = "%s"
    }

    stack_configuration = {
        name_pattern = "{stack_name}"
        iac_type = "%s"

        vcs_info_with_patterns = {
            provider_id = "%s"
            repo_name = "%s"
            path_pattern = "{stack_path}"
        }
    }

    substitute_parameters = [
		{
			key = "stack_name"
			description = "any name you want"
		},
		{
			key = "stack_path"
			description = "path"
		}
	]
}
`, tfBlueprintResourceResource, blueprintTfResourceName, blueprintName, blueprintDescription, blueprintProviderId,
					blueprintRepoName, blueprintRepoPath, blueprintStackConfigurationIacType, blueprintProviderId, blueprintRepoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(blueprintResourceName(blueprintTfResourceName), "id"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "name", blueprintName),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "description", blueprintDescription),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "blueprint_vcs_info.provider_id", blueprintProviderId),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "blueprint_vcs_info.repo_name", blueprintRepoName),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "blueprint_vcs_info.path", blueprintRepoPath),
					resource.TestCheckNoResourceAttr(blueprintResourceName(blueprintTfResourceName), "blueprint_vcs_info.branch"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.name_pattern", "{stack_name}"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.iac_type", blueprintStackConfigurationIacType),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.vcs_info_with_patterns.provider_id", blueprintProviderId),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.vcs_info_with_patterns.repo_name", blueprintRepoName),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.vcs_info_with_patterns.path_pattern", "{stack_path}"),
					resource.TestCheckNoResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.vcs_info_with_patterns.branch_pattern"),
					resource.TestCheckNoResourceAttr(blueprintResourceName(blueprintTfResourceName), "stack_configuration.deployment_approval_policy"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.#", "2"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.0.key", "stack_name"),
					resource.TestCheckResourceAttrSet(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.0.description"),
					resource.TestCheckResourceAttr(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.1.key", "stack_path"),
					resource.TestCheckResourceAttrSet(blueprintResourceName(blueprintTfResourceName), "substitute_parameters.1.description"),
					resource.TestCheckNoResourceAttr(blueprintResourceName(blueprintTfResourceName), "skip_plan_on_stack_initialization"),
					resource.TestCheckNoResourceAttr(blueprintResourceName(blueprintTfResourceName), "auto_approve_apply_on_initialization"),
				),
			},
			// validate no drift step
			test_helpers.GetValidateNoDriftStep(),
		},
	})
}

func blueprintResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfBlueprintResourceResource, s)
}
