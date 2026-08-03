resource "cm_custom_role" "platform_admin" {
  name        = "Platform Admin"
  description = "Organization level role that grants stack visibility and creation across the organization"
  type        = "organizationRole"

  permissions = [
    {
      # Organization roles use the plural `names`. The singular `name` is for namespace roles only.
      # `org:stack:create` requires `org:stack:read` to be granted in the same block.
      names = ["org:stack:read", "org:namespace:read", "org:stack:create"]
    },
    {
      # Restrictions are only accepted on permissions that support them, such as `cloudAccount:*`.
      # ControlMonkey rejects restrictions on `org:*` actions.
      names = ["cloudAccount:insights"]
      restrictions = [
        {
          # Fields within a single restriction are combined with AND.
          cloud_provider   = "aws"
          cloud_account_id = "123456789012"
        },
        {
          # Multiple restrictions are combined with OR, so this one widens the permission.
          cloud_provider = "azure"
        }
      ]
    }
  ]
}
