resource "cm_namespace" "dev_cross_account_namespace" {
  name        = "Dev"
  description = "AWS dev env that consists of dev-main, dev-infra & dev-monitoring AWS accounts"

  external_credentials = [
    {
      external_credentials_id = "ext-aws-dev-main" # default credentials
      type                    = "awsAssumeRole"
    },
    {
      external_credentials_id = "ext-aws-dev-infra"
      type                    = "awsAssumeRole"
      aws_profile_name        = "dev-infra"
    },
    {
      external_credentials_id = "ext-aws-dev-monitoring"
      type                    = "awsAssumeRole"
      aws_profile_name        = "dev-monitoring"
    }
  ]
}
