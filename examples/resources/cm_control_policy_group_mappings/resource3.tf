resource "cm_control_policy_group_mappings" "policy_group" {
  control_policy_group_id = cm_control_policy_group.policy_group.id

  targets = [
    {
      target_id             = cm_namespace.namespace.id
      target_type           = "namespace"
      enforcement_level     = "bySeverity"
      override_enforcements = [
        {
          control_policy_id = cm_control_policy.control_policy.id
          enforcement_level = "softMandatory"
        },
      ]
    },
  ]
}