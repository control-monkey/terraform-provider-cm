resource "cm_variable" "default_log" {
  scope           = "namespace"
  scope_id        = cm_namespace.namespace.id
  key             = "TF_LOG"
  type            = "envVar"
  value           = "ERROR"
  is_sensitive    = false
  is_overridable  = true
}
