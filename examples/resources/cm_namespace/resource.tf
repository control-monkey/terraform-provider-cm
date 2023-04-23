resource "cm_namespace" "namespace" {
  name        = "namespace1"
  description = "first namespace"

  external_credentials = [
    {
      type                    = "awsAssumeRole"
      external_credentials_id = "ext-123"
    },
    {
      type                    = "awsAssumeRole"
      external_credentials_id = "ext-1234"
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
