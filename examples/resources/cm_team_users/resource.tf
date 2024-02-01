resource "cm_team_users" "dev_team_users" {
  team_id = cm_team.dev_team.id

  users = [
    {
      email = "example1@email.com"
    },
    {
      email = "example2@email.com"
    }
  ]
}

resource "cm_team" "dev_team" {
  name = "Dev Team"
}

