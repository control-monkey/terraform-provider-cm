resource "cm_control_policy_group" "policy_group" {
  name             = "Mandatory policies"
  description      = "These policies should be applied to all resources"
  control_policies = [
    {
      control_policy_id = cm_control_policy.control_policy_1.id
      severity          = "high"
    },
    {
      control_policy_id = cm_control_policy.control_policy_2.id
      severity          = "medium"
    },
  ]
}