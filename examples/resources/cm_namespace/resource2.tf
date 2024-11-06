resource "cm_namespace" "dev_cross_account_namespace" {
  name        = "Dev"
  description = "AWS dev env that consists of dev-main, dev-infra & dev-monitoring AWS accounts"

  external_credentials = [
    {
      external_credentials_id = "ext-123"
      type                    = "awsAssumeRole"
    },
    {
      external_credentials_id = "ext-456"
      type                    = "awsAssumeRole"
      aws_profile_name        = "dev-infra"
    },
    {
      external_credentials_id = "ext-789"
      type                    = "awsAssumeRole"
      aws_profile_name        = "dev-monitoring"
    }
  ]
}
