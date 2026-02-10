data "cm_external_credential" "aws_dev_creds" {
  name   = "dev-account-credential"
  vendor = "aws"
}

resource "cm_namespace" "dev_namespace" {
  name        = "Dev"
  description = "AWS dev env"

  external_credentials = [
    {
      type                    = "awsAssumeRole"
      external_credentials_id = data.cm_external_credential.aws_dev_creds.id
    }
  ]
}