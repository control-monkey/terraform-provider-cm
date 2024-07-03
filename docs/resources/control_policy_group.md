---
page_title: "cm_control_policy_group Resource - terraform-provider-cm"
subcategory: ""
description: |-
  Creates, updates and destroys control policy groups.
---

# cm_control_policy_group (Resource)

Creates, updates and destroys control policy groups.

## Example Usage
```terraform
resource "cm_control_policy_group" "policy_group" {
  name             = "Mandatory policies"
  description      = "These policies should be applied to all resources"
  control_policies = [
    {
      control_policy_id = cm_control_policy.control_policy_1.id
      severity          = "high"
    },
    {
      control_policy_id = cm_control_policy.control_policy_2.id
      severity          = "medium"
    },
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `control_policies` (Attributes List) List of control policies to enforce. (see [below for nested schema](#nestedatt--control_policies))
- `name` (String) The name of the control policy group.

### Optional

- `description` (String) The description of the control policy group.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--control_policies"></a>
### Nested Schema for `control_policies`

Required:

- `control_policy_id` (String) The ControlMonkey unique ID of the control policy.

Optional:

- `severity` (String) The severity of the control policy within the group is determined by the severity parameter. This parameter becomes effective only when a mapping is established in [cm_control_policy_group_mappings](https://registry.terraform.io/providers/control-monkey/cm/latest/docs/resources/control_policy_group_mappings) and the enforcementLevel is set to 'bySeverity'. Allowed values: [low, medium, high, critical].

## Import

`cm_control_policy_group` can be imported using the ID of the Control Policy Group, e.g.

```shell
terraform import cm_control_policy_group.policy_group polg-123
```