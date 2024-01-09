resource "cm_variable" "default_log" {
  scope           = "namespace"
  scope_id        = "ns-stage"
  key             = "TF_LOG"
  type            = "envVar"
  value           = "ERROR"
  is_sensitive    = false
  is_overridable  = true
}
