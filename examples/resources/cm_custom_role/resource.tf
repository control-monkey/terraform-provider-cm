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