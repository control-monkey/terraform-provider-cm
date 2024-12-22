resource "cm_disaster_recovery_configuration" "autoManagedStackConfiguration" {
  scope            = "aws"
  cloud_account_id = "123456789"

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
            "path": "s3/us-east-1"
          }

          "awsQuery": {
            "region": "us-east-1",
            "services": ["AWS::S3"],
            "resourceTypes": ["AWS::S3::Bucket"]
          }
        },
        {
          "vcsInfo": {
            "path": "ecs/us-east-1"
          }

          "awsQuery": {
            "region": "us-east-1",
            "services": ["AWS::ECS"]
          }
        },
        {
          "vcsInfo": {
            "path": "ec2/us-east-1"
          }

          "awsQuery": {
            "region": "us-east-1",
            "services": ["AWS::EC2"]
          }
        }
      ]
    )
  }
}

