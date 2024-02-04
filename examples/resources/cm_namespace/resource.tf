resource "cm_namespace" "dev_namespace" {
  name        = "Dev"
  description = "AWS dev env"

  external_credentials = [
    {
      type                    = "awsAssumeRole"
      external_credentials_id = "ext-123"
    }
  ]
}
