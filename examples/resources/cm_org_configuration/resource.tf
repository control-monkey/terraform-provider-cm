resource "cm_org_configuration" "org_configuration" {
  iac_config = {
    terraform_version = "1.5.5"
  }

  s3_state_files_locations = [
    {
      bucket_name    = "bucket-example"
      bucket_region  = "us-east-1"
      aws_account_id = "123456789"
    },
  ]

  suppressed_resources = {
    managed_by_tags = [
      {
        key = "aws:eks:cluster-name"
      },
    ]
  }

  report_configurations = [
    {
      enabled    = true
      type       = "weeklyReport"
      recipients = {
        all_admins = true
      },
    },
  ]
}