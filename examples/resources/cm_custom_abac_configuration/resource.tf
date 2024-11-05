resource "cm_team" "team" {
  name = "SRE Team"
}

resource "cm_custom_abac_configuration" "abac_configuration" {
  custom_abac_id = "xxxx"
  name           = "SRE Team ABAC Configuration"
  roles = [
    {
      org_id   = "o-123"
      org_role = "admin"
    },
    {
      org_id   = "o-123"
      org_role = "member"
      team_ids = [cm_team.team.id]
    }
  ]
}
