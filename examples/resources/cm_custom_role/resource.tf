# A namespace role. `type` is optional and defaults to namespaceRole, but setting it explicitly
# makes clear which attributes are legal: namespace roles use the singular `permissions.name` and
# may set `stack_restriction`.
resource "cm_custom_role" "custom_role" {
  name        = "Create Stack Role"
  description = "This role allows users to create stack and launch a stack from an ephemeral template"
  type        = "namespaceRole"
  permissions = [
    {
      name = "stack:create"
    },
    {
      name = "stack:createFromTemplate"
    }
  ]
}
