resource "cm_namespace_permissions" "stage_namespace_permissions" {
  namespace_id = cm_namespace.stage_namespace.id

  permissions = [
    {
      user_email = "example@email.com"
      role       = "viewer"
    },
    {
      team_id = cm_team.stage_team_developers.id
      role    = "deployer"
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

resource "cm_namespace" "stage_namespace" {
  name = "Stage"
}

resource "cm_team" "stage_team_developers" {
  name = "Stage Team Developers"
}

resource "cm_team" "stage_team_it" {
  name = "Stage Team IT"
}

