---
page_title: "cm_template Resource - terraform-provider-cm"
subcategory: ""
description: |-
  Creates, updates and destroys templates for ephemeral stack.
---

# cm_template (Resource)

Creates, updates and destroys templates for ephemeral stack.

## Example Usage

### Basic template for an ephemeral stack that provides self-service capabilities for developers.
```terraform
resource "cm_template" "template_developers" {
  name = "Dev Self-Service Template"
  iac_type = "terraform"
  description = "Self service on Dev environment for developers"

  vcs_info = {
    provider_id = "vcsp-123"
    repo_name = "terraform"
    path = "dev/self-service"
  }
}
```

### Template for provisioning an ephemeral demo environment stack with default TTL and maximum TTL limit.
```terraform
resource "cm_template" "temporary_demo_template" {
  name = "Demo Template"
  iac_type = "terraform"
  description = "Template for temporary demo environment with TTL"

  vcs_info = {
    provider_id = "vcsp-123"
    repo_name = "terraform"
    path = "demo/template"
    branch = "demo"
  }

  policy = {
    ttl_config = {
      max_ttl = {
        type  = "days"
        value = "1"
      }
      default_ttl = {
        type  = "hours"
        value = "5"
      }
    }
  }
}
```

### Template for provisioning an ephemeral RDS stack with TTL, including the required variables that must be provided by the provisioner of the ephemeral stack.
```terraform
resource "cm_template" "rds_template" {
  name     = "Ephemeral RDS For R&D"
  iac_type = "terraform"

  vcs_info = {
    provider_id = "vcsp-123"
    repo_name   = "terraform"
    path        = "templates/rds"
  }

  policy = {
    ttl_config = {
      max_ttl = {
        type  = "days"
        value = "10"
      }
      default_ttl = {
        type  = "days"
        value = "5"
      }
    }
  }
}

resource "cm_variable" "creator_variable" {
  scope          = "template"
  scopeId        = cm_template.rds_template.id
  key            = "creator_name"
  type           = "tfVar"
  is_sensitive   = false
  is_overridable = true
  is_required    = true
}

resource "cm_variable" "allowed_instances_variable" {
  scope          = "template"
  scopeId        = cm_template.rds_template.id
  key            = "instance_type"
  type           = "tfVar"
  is_sensitive   = false
  is_overridable = true
  is_required    = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `iac_type` (String) IaC type of the template. Allowed values: [terraform, terragrunt, opentofu].
- `name` (String) The name of the template.
- `vcs_info` (Attributes) The configuration of the version control to which the template is attached. (see [below for nested schema](#nestedatt--vcs_info))

### Optional

- `description` (String) The description of the template.
- `policy` (Attributes) The policy of the template. (see [below for nested schema](#nestedatt--policy))
- `skip_state_refresh_on_destroy` (Boolean) When enabled, the state will not get refreshed before planning the destroy operation.

### Read-Only

- `id` (String) The unique ID of the template.

<a id="nestedatt--vcs_info"></a>
### Nested Schema for `vcs_info`

Required:

- `provider_id` (String) The ControlMonkey unique ID of the connected version control system.
- `repo_name` (String) The name of the version control repository.

Optional:

- `branch` (String) The branch that triggers the deployment of the ephemeral stack from the template. If no branch is specified, the default branch of the repository will be used.
- `path` (String) The path to a chosen directory from the root. Default path is root directory


<a id="nestedatt--policy"></a>
### Nested Schema for `policy`

Optional:

- `ttl_config` (Attributes) The time to live config of the template policy. (see [below for nested schema](#nestedatt--policy--ttl_config))

<a id="nestedatt--policy--ttl_config"></a>
### Nested Schema for `policy.ttl_config`

Required:

- `default_ttl` (Attributes) (see [below for nested schema](#nestedatt--policy--ttl_config--default_ttl))
- `max_ttl` (Attributes) (see [below for nested schema](#nestedatt--policy--ttl_config--max_ttl))

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

`cm_template` can be imported using the ID of the Template for ephemeral stack, e.g.

```shell
terraform import cm_template.template tmpl-123
```