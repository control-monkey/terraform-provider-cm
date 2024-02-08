resource "cm_control_policy_mappings" "no_public_bucket_prod" {
  control_policy_id = "pol-123"

  targets = [
    {
      target_id         = cm_namespace.namespace.id
      target_type       = "namespace"
      enforcement_level = "hardMandatory"
    },
    {
      target_id         = cm_stack.stack.id
      target_type       = "stack"
      enforcement_level = "softMandatory"
    },
  ]
}