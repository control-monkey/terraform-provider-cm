resource "cm_variable" "var_stack" {
  scope          = "stack"
  scope_id       = "stk-nb294bzary"
  key            = "TfProvider2"
  type           = "tfVar"
  value          = "Value1"
  is_sensitive   = false
  is_overridable = false
}

resource "cm_variable" "var_namespace" {
  scope          = "namespace"
  scope_id       = "ns-1cizwgs7jv"
  key            = "TfProvider"
  type           = "tfVar"
  value          = "Value2"
  is_sensitive   = false
  is_overridable = false
}

resource "cm_variable" "var_template" {
  scope          = "template"
  scope_id       = "tmpl-0uzcrl6ytn"
  key            = "Value_template1"
  type           = "tfVar"
  value          = "Value3"
  is_sensitive   = false
  is_overridable = true
  is_required    = true
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