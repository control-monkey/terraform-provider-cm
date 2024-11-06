resource "cm_namespace" "prod_namespace" {
  name = "Prod"

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
