resource "cm_disaster_recovery_configuration" "dr_config" {
  scope            = "aws"
  cloud_account_id = "123465789"

  backup_strategy = {
    include_managed_resources = true
    mode                      = "manual"

    vcs_info = {
      provider_id = "vcsp-123"
      repo_name   = "terraform-backup"
      branch      = "main"
    }

    groups_json = jsonencode(
      [
        {
          "vcsInfo": {
            "path": "ec2/instances/us-east-1"
          }

          "awsQuery": {
            "region": "us-east-1",
            "services": ["AWS::EC2"],
            "resourceTypes": ["AWS::EC2::Instance"]

            "tags": [
              {
                "key": "Env",
                "value": "Prod"
              }
            ]
          }
        },
      ]
    )
  }
}