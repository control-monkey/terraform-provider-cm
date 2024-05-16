data "cm_namespace" "dev_namespace" {
  name = "Dev"
}


resource "cm_namespace_permissions" "dev_namespace_permissions" {
  namespace_id = data.cm_namespace.dev_namespace.id

  permissions = [
    {
      user_email = "example@email.com"
      role       = "viewer"
    }
  ]
}
