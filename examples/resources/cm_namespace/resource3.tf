resource "cm_namespace" "dev_namespace" {
  name        = "Dev"
  description = "AWS dev env"

  external_credentials = [
    {
      type                    = "awsAssumeRole"
      external_credentials_id = "ext-123"
    }
  ]

  policy = {
    ttl_config = {
      max_ttl = {
        type  = "days"
        value = "2"
      }
      default_ttl = {
        type  = "hours"
        value = "3"
      }
    }
  }
}
