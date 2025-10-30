---
page_title: "cm_blueprint Resource - terraform-provider-cm"
subcategory: ""
description: |-
  Creates, updates and destroys blueprints. For more information: ControlMonkey Documentation https://docs.controlmonkey.io/main-concepts/self-service-templates/persistent-template
---

# cm_blueprint (Resource)

Creates, updates and destroys blueprints. For more information: [ControlMonkey Documentation](https://docs.controlmonkey.io/main-concepts/self-service-templates/persistent-template)

## Learn More

- [Self-service Infrastructure](https://controlmonkey.io/blog/self-service-infrastructure/)
- [Self-service templates support for Terragrunt & OpenTofu](https://controlmonkey.io/news/self-service-templates-support-for-terragrunt-opentofu/)
- [Variable Conditions for Self-service Infrastructure](https://controlmonkey.io/news/variable-conditions-for-self-service-infrastructure/)

## Example Usage
~> **NOTE:** All dynamic parameters referenced in `name_pattern`, `path_pattern`, and `branch_pattern` must be explicitly defined in `substitute_parameters`.

```terraform
resource "cm_blueprint" "blueprint" {
  name = "Blueprint example"

  blueprint_vcs_info = {
    provider_id = "vcsp-github"
    repo_name   = "terraform"
    path        = "infra/blueprints/vpc"
  }

  stack_configuration = {
    name_pattern = "{env}-{region}-{service}"
    iac_type     = "terraform"

    vcs_info_with_patterns = {
      provider_id  = "vcsp-github"
      repo_name    = "terraform"
      path_pattern = "infra/{env}/{region}/{service}"
    }
  }

  substitute_parameters = [
    {
      key         = "env"
      description = "The name of the environment e.g prod, stage"
    },
    {
      key         = "region"
      description = "The region in the cloud e.g us-east-1"
    },
    {
      key         = "service"
      description = "The name of the service e.g ec2, s3"
    }
  ]
}
```

### This example demonstrates a blueprint stored in the terraform repository, while stacks launched from it are pushed to a separate repository named terraform-developers. Deployment of the new infrastructure requires approval from two distinct users. Additionally, the substitute parameter env includes specific conditions for its value, which must be provided by the user when launching the stack.
```terraform
resource "cm_blueprint" "blueprint" {
  name = "Blueprint example"

  blueprint_vcs_info = {
    provider_id = "vcsp-github"
    repo_name   = "terraform"
    path        = "infra/blueprints/vpc"
  }

  stack_configuration = {
    name_pattern = "{env}-{region}-{service}"
    iac_type     = "terraform"

    vcs_info_with_patterns = {
      provider_id    = "vcsp-github"
      repo_name      = "terraform-developers"
      path_pattern   = "infra/{env}/{region}/{service}"
      branch_pattern = "{env}-development"
    }

    deployment_approval_policy = {
      rules = [
        {
          type = "requireTwoApprovals"
        }
      ]
    }
  }

  substitute_parameters = [
    {
      key         = "env"
      description = "The name of the environment in lower case e.g prod, stage"
      value_conditions = [
        {
          operator = "in"
          values = ["prod", "stage", "dev"]
        }
      ]
    },
    {
      key         = "region"
      description = "The region in the cloud e.g us-east-1"
    },
    {
      key         = "service"
      description = "The name of the service in lower case e.g ec2, s3"
    }
  ]
}
```

### This example features a ready-to-run blueprint that is already mapped to a namespace where new stacks can be launched. Deployment of a new stack requires team approval. Additionally, the code_dynamic_parameter substitute parameter is used in the blueprint files, and its value—provided by the user during stack launch—will replace all occurrences of the parameter in the code.
```terraform
data "cm_namespace" "prod_namespace" {
  name = "Prod"
}

data "cm_team" "prod_devops_team" {
  name = "Prod DevOps Team"
}

resource "cm_blueprint_namespace_mappings" "vpc_blueprint_namespace_mappings" {
  blueprint_id = cm_blueprint.blueprint.id

  namespaces = [
    {
      namespace_id = data.cm_namespace.prod_namespace.id
    }
  ]
}

resource "cm_blueprint" "blueprint" {
  name = "Blueprint example"

  blueprint_vcs_info = {
    provider_id = "vcsp-github"
    repo_name   = "terraform"
    path        = "infra/blueprints/vpc"
    branch      = "cm-blueprints"
  }

  stack_configuration = {
    name_pattern = "{env}-{region}-{service}"
    iac_type     = "terraform"

    vcs_info_with_patterns = {
      provider_id    = "vcsp-github"
      repo_name      = "terraform"
      path_pattern   = "infra/{env}/{region}/{service}"
      branch_pattern = "{env}-main"
    }

    deployment_approval_policy = {
      rules = [
        {
          type = "requireTeamsApproval"
          parameters = jsonencode({
            teams = [data.cm_team.prod_devops_team.id]
          })
        },
      ]
    }
  }

  substitute_parameters = [
    {
      key         = "env"
      description = "The name of the environment in lower case e.g prod, stage"
      value_conditions = [
        {
          operator = "in"
          values = ["prod", "stage", "dev"]
        }
      ]
    },
    {
      key         = "region"
      description = "The region in the cloud e.g us-east-1"
    },
    {
      key         = "service"
      description = "The name of the service in lower case e.g ec2, s3"
    },
    {
      key         = "code_dynamic_parameter"
      description = "The provided value will replace all occurrences of {code_dynamic_parameter} in the code"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `blueprint_vcs_info` (Attributes) Configuration details for the version control system storing the blueprint. (see [below for nested schema](#nestedatt--blueprint_vcs_info))
- `name` (String) The name of the blueprint.
- `stack_configuration` (Attributes) The configuration for creating new persistent stacks from the blueprint. (see [below for nested schema](#nestedatt--stack_configuration))
- `substitute_parameters` (Attributes List) Define dynamic placeholders (`{parameter_name}`) used in patterns (e.g., `name_pattern`, `path_pattern`) or Terraform files. Users will supply values for these parameters when launching stacks. (see [below for nested schema](#nestedatt--substitute_parameters))

### Optional

- `auto_approve_apply_on_initialization` (Boolean) If enabled (`true`), the stack’s initial deployment will automatically apply changes after the pull request is merged, bypassing manual approval.
- `description` (String) The description of the blueprint.
- `policy` (Attributes) The policy of the blueprint. (see [below for nested schema](#nestedatt--policy))
- `skip_plan_on_stack_initialization` (Boolean) If enabled (`true`), an automatic plan will not be triggered on the initial pull request.

### Read-Only

- `id` (String) The ID of the blueprint.

<a id="nestedatt--blueprint_vcs_info"></a>
### Nested Schema for `blueprint_vcs_info`

Required:

- `path` (String) The relative path to the directory containing the blueprint files, starting from the root of the repository.
- `provider_id` (String) The ControlMonkey unique ID of the connected version control system.
- `repo_name` (String) The name of the version control repository.

Optional:

- `branch` (String) The branch in which the blueprint is located. When no branch is given, the default branch of the repository is chosen.


<a id="nestedatt--stack_configuration"></a>
### Nested Schema for `stack_configuration`

Required:

- `iac_type` (String) IaC type of the template. Allowed values: [terraform, terragrunt, opentofu].
- `name_pattern` (String) A pattern used to name persistent stacks created from the blueprint. The pattern must include at least one dynamic substitute parameter (e.g., `{region}-{service}`).
- `vcs_info_with_patterns` (Attributes) Configuration details for the version control system where the stack files generated from the blueprint will be stored. (see [below for nested schema](#nestedatt--stack_configuration--vcs_info_with_patterns))

Optional:

- `auto_sync` (Attributes) Set up auto sync configurations. (see [below for nested schema](#nestedatt--stack_configuration--auto_sync))
- `deployment_approval_policy` (Attributes) Set up requirements to approve a deployment (see [below for nested schema](#nestedatt--stack_configuration--deployment_approval_policy))
- `iac_config` (Attributes) IaC configuration. (see [below for nested schema](#nestedatt--stack_configuration--iac_config))
- `run_trigger` (Attributes) Glob patterns to specify additional paths that should trigger a stack run. (see [below for nested schema](#nestedatt--stack_configuration--run_trigger))

<a id="nestedatt--stack_configuration--vcs_info_with_patterns"></a>
### Nested Schema for `stack_configuration.vcs_info_with_patterns`

Required:

- `path_pattern` (String) A pattern to a new path in the repository to which new persistent stack files created from the blueprint will be pushed. This field requires at least one substitute parameter.
- `provider_id` (String) The ControlMonkey unique ID of the connected version control system.
- `repo_name` (String) The name of the version control repository.

Optional:

- `branch_pattern` (String) The target branch for new pull requests containing the new stack files. Substitute parameters (e.g., `{branch}-{env}`) are supported.


<a id="nestedatt--stack_configuration--auto_sync"></a>
### Nested Schema for `stack_configuration.auto_sync`

Optional:

- `deploy_when_drift_detected` (Boolean) If set to `true`, a deployment will start automatically upon detecting a drift or multiple drifts


<a id="nestedatt--stack_configuration--deployment_approval_policy"></a>
### Nested Schema for `stack_configuration.deployment_approval_policy`

Required:

- `rules` (Attributes List) Set up rules for approving deployment processes. At least one rule should be configured (see [below for nested schema](#nestedatt--stack_configuration--deployment_approval_policy--rules))

<a id="nestedatt--stack_configuration--deployment_approval_policy--rules"></a>
### Nested Schema for `stack_configuration.deployment_approval_policy.rules`

Required:

- `type` (String) The type of the rule. Find supported types [here](https://docs.controlmonkey.io/controlmonkey-api/api-enumerations#deployment-approval-policy-rule-types)

Optional:

- `parameters` (String) JSON format of the rule parameters according to the `type`. Find supported parameters [here](https://docs.controlmonkey.io/controlmonkey-api/approval-policy-rules)



<a id="nestedatt--stack_configuration--iac_config"></a>
### Nested Schema for `stack_configuration.iac_config`

Optional:

- `is_terragrunt_run_all` (Boolean) When using terragrunt, as long as this field is set to `True`, this field will execute "run-all" commands on multiple modules for init/plan/apply
- `opentofu_version` (String) the OpenTofu version that will be used for tofu operations.
- `terraform_version` (String) the Terraform version that will be used for terraform operations.
- `terragrunt_version` (String) the Terragrunt version that will be used for terragrunt operations.
- `var_files` (List of String) Custom variable files to pass on to Terraform. For more information: [ControlMonkey Docs](https://docs.controlmonkey.io/main-concepts/stack/stack-settings#var-files)


<a id="nestedatt--stack_configuration--run_trigger"></a>
### Nested Schema for `stack_configuration.run_trigger`

Optional:

- `exclude_patterns` (List of String) Patterns that will not trigger a stack run.
- `patterns` (List of String) Patterns that trigger a stack run.



<a id="nestedatt--substitute_parameters"></a>
### Nested Schema for `substitute_parameters`

Required:

- `description` (String) A description of the parameter. Users launching stacks from this blueprint will reference this description to assign values. Providing a clear, meaningful description is highly recommended.
- `key` (String) The key for the substitute parameter excluding the curly braces. For example, if the Terraform file contains `{replace-me}`, the key should be `replace-me`.

Optional:

- `value_conditions` (Attributes List) Specify conditions for the variable value using an operator and another value. Typically used for stacks launched from templates. For more information: [ControlMonkey Docs] (https://docs.controlmonkey.io/main-concepts/variables/variable-conditions) (see [below for nested schema](#nestedatt--substitute_parameters--value_conditions))

<a id="nestedatt--substitute_parameters--value_conditions"></a>
### Nested Schema for `substitute_parameters.value_conditions`

Required:

- `operator` (String) Logical operators. Allowed values: [ne, gt, gte, lt, lte, in, startsWith, contains].

Optional:

- `value` (String) The value associated with the operator. Input a number or string depending on the chosen operator. Use `values` field for operator of type `in`
- `values` (List of String) A list of strings when using operator type `in`. For other operators use `value`



<a id="nestedatt--policy"></a>
### Nested Schema for `policy`

Optional:

- `ttl_config` (Attributes) The time to live config of the blueprint policy. (see [below for nested schema](#nestedatt--policy--ttl_config))

<a id="nestedatt--policy--ttl_config"></a>
### Nested Schema for `policy.ttl_config`

Required:

- `default_ttl` (Attributes) The default time to live configuration for the blueprint. (see [below for nested schema](#nestedatt--policy--ttl_config--default_ttl))
- `max_ttl` (Attributes) The maximum time to live configuration for the blueprint. (see [below for nested schema](#nestedatt--policy--ttl_config--max_ttl))

Optional:

- `open_cleanup_pr_on_ttl_termination` (Boolean) When enabled, a PR will automatically open to remove the stack directory from the repository after the stack is terminated due to TTL expiration.

<a id="nestedatt--policy--ttl_config--default_ttl"></a>
### Nested Schema for `policy.ttl_config.default_ttl`

Required:

- `type` (String) The type of the ttl. Allowed values: [hours, days].
- `value` (Number) The value that corresponds the type


<a id="nestedatt--policy--ttl_config--max_ttl"></a>
### Nested Schema for `policy.ttl_config.max_ttl`

Required:

- `type` (String) The type of the ttl. Allowed values: [hours, days].
- `value` (Number) The value that corresponds the type

## Import

`cm_blueprint` can be imported using the ID of the Blueprint, e.g.

```shell
terraform import cm_blueprint.blueprint blp-123
```