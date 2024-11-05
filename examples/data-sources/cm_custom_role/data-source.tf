data "cm_custom_role" "create_stack_role" {
  name = "Create Stack Role"
}

data "cm_namespace" "dev_namespace" {
  name = "Dev"
}

data "cm_team" "developers_team" {
  name = "Developers"
}


resource "cm_namespace_permissions" "dev_namespace_permissions" {
  namespace_id = data.cm_namespace.dev_namespace.id

  permissions = [
    {
      team_id = data.cm_team.developers_team.id
      custom_role_id = data.cm_custom_role.create_stack_role.id
    }
  ]
}
