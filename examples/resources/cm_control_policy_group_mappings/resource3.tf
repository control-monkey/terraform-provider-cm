resource "cm_control_policy_group_mappings" "policy_group" {
  control_policy_group_id = "cmpg-123"

  targets = [
    {
      target_id         = cm_namespace.namespace.id
      target_type       = "namespace"
      enforcement_level = "softMandatory"
      override_enforcements = [
        {
          control_policy_id = "cmp-123"
          enforcement_level = "warning"
        },
        {
          control_policy_id = "cmp-124"
          enforcement_level = "hardMandatory"
          stack_ids = [cm_stack.stack2.id]
        }
      ]
    },
  ]
}