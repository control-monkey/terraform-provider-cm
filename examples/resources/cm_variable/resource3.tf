resource "cm_variable" "stage_ephemeral_stacks_default_volume_size" {
  scope            = "template"
  scope_id         = cm_template.template.id
  key              = "volume_size"
  type             = "tfVar"
  value            = 8
  description      = "Default volume size for ephemeral stacks in GB. Can be overridden up to 50GB"
  is_sensitive     = false
  is_overridable   = true
  value_conditions = [
    {
      operator = "lte"
      value    = 50
    }
  ]
}
