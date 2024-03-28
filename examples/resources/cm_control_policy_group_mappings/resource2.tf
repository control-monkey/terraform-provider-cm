resource "cm_control_policy_group" "policy_group" {
  name             = "Mandatory policies"
  description      = "These policies should be applied to all resources"
  control_policies = [
    {
      control_policy_id = cm_control_policy.control_policy_1.id
      severity          = "medium" // equals to enforcement level warning
    },
    {
      control_policy_id = cm_control_policy.control_policy_2.id
      severity          = "critical" // equals to enforcement level hardMandatory
    },
  ]
}

resource "cm_control_policy_group_mappings" "policy_group_mappings" {
  control_policy_group_id = cm_control_policy_group.policy_group.id

  targets = [
    {
      target_id         = cm_namespace.namespace.id
      target_type       = "namespace"
      enforcement_level = "bySeverity"
    },
  ]
}