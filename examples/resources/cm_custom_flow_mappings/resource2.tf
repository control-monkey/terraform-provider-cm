# Map a custom flow to every namespace in the organization
resource "cm_custom_flow_mappings" "all_namespaces" {
  custom_flow_id = cm_custom_flow.opa_validation.id

  targets = [
    {
      target_id   = "ALL"
      target_type = "namespace"
    },
  ]
}
