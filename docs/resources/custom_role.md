---
page_title: "cm_custom_role Resource - terraform-provider-cm"
subcategory: ""
description: |-
  Creates, updates and destroys custom roles.
---

# cm_custom_role (Resource)

Creates, updates and destroys custom roles.

## Example Usage
```terraform
resource "cm_custom_role" "custom_role" {
  name        = "Create Stack Role"
  description = "This role allows users to create stack and launch a stack from an ephemeral template"
  permissions = [
    {
      name = "stack:create"
    },
    {
      name = "stack:createFromTemplate"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the role.

### Optional

- `description` (String) The description of the role.
- `permissions` (Attributes List) List of permissions allowed by the role. (see [below for nested schema](#nestedatt--permissions))
- `stack_restriction` (String) Restrict stack operations with supported types. Learn more [here](https://docs.controlmonkey.io/administration/users-and-roles/custom-roles). Find supported types [here](https://docs.controlmonkey.io/controlmonkey-api/api-enumerations#stack-restriction-types).

### Read-Only

- `id` (String) The ID of the custom role.

<a id="nestedatt--permissions"></a>
### Nested Schema for `permissions`

Required:

- `name` (String) The type of the permission. Find supported types [here](https://docs.controlmonkey.io/controlmonkey-api/api-enumerations#custom-role-permission-types).

## Import

`cm_custom_role` can be imported using the ID of the Custom Role, e.g.

```shell
terraform import cm_custom_role.custom_role cro-123
```