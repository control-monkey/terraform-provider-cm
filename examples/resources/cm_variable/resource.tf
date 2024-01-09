resource "cm_variable" "organization_mandatory_variable" {
  scope          = "organization"
  key            = "ORG_BU_IDENTIFIER"
  type           = "envVar"
  value          = "162"
  description    = "Environment variable across all business units"
  is_sensitive   = false
  is_overridable = false
}