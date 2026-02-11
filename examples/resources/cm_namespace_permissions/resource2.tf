resource "cm_namespace_permissions" "staging_stack_permissions" {
  stack_id = cm_stack.staging_stack_permissions.id

  permissions = [
    {
      user_email = "example@email.com"
      role       = "viewer"
    },
    {
      team_id = cm_team.stage_team_developers.id
      role    = "viewer"
    },
    {
      programmatic_username = "automation-user"
      role                  = "admin"
    },
    {
      team_id        = cm_team.stage_team_it.id
      custom_role_id = "cro-123"
    },
  ]
}
