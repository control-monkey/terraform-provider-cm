resource "cm_custom_flow_mappings" "opa_validation_mappings" {
  custom_flow_id = cm_custom_flow.opa_validation.id

  targets = [
    {
      target_id   = cm_namespace.namespace.id
      target_type = "namespace"
    },
    {
      target_id   = cm_stack.stack.id
      target_type = "stack"
    },
  ]
}
