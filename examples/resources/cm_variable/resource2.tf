resource "cm_variable" "dev_basic_app_mandatory_owner" {
  scope          = "template"
  scope_id       = cm_template.template.id
  key            = "Owner"
  type           = "tfVar"
  description    = "Provide your name. It will be used as the tag value of the the same key"
  is_sensitive   = false
  is_overridable = true
  is_required    = true
}