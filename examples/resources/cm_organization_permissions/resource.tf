resource "cm_custom_role" "platform_admin" {
  name = "platform admin"
  type = "organizationRole"

  permissions = [
    {
      names = ["org:stack:read", "org:stack:create"]
    }
  ]
}

resource "cm_organization_permissions" "platform_team" {
  team_id = cm_team.platform.id

  permissions = [
    {
      custom_role_id = cm_custom_role.platform_admin.id
    }
  ]
}
