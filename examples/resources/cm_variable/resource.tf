resource "cm_variable" "var_namespace" {
  scope          = "namespace"
  scope_id       = "ns-123"
  key            = "TfProvider"
  type           = "tfVar"
  value          = "Value2"
  is_sensitive   = false
  is_overridable = false
}

resource "cm_variable" "var_org" {
  scope          = "organization"
  key            = "orgKey"
  type           = "envVar"
  value          = "Value4"
  is_sensitive   = false
  is_overridable = false
  is_required    = false
}