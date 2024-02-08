resource "cm_control_policy_group_mappings" "policy_group" {
  control_policy_group_id = "cmpg-123"

  targets = [
    {
      target_id         = cm_namespace.namespace.id
      target_type       = "namespace"
      enforcement_level = "bySeverity"
    },
    {
      target_id         = cm_stack.stack.id
      target_type       = "stack"
      enforcement_level = "softMandatory"
    },
  ]
}