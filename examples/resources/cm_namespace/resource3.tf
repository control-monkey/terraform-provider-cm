data "cm_team" "team_devops" {
  name = "DevOps Team"
}

data "cm_team" "team_prod" {
  name = "Prod Team"
}

resource "cm_namespace" "prod_namespace" {
  name = "Prod"

  external_credentials = [
    {
      type                    = "awsAssumeRole"
      external_credentials_id = "ext-123"
    }
  ]

  deployment_approval_policy = {
    override_behavior = "deny"
    rules = [
      {
        type = "requireTeamsApproval"
        parameters = jsonencode({
          teams = [cm_team.team_devops.id, cm_team.team_prod.id]
        })
      },
    ]
  }
}
