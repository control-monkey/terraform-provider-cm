data "cm_control_policy" "cost_threshold_to_approve" {
  name = "Cost Threshold"
}


resource "cm_control_policy_group" "cost_policy_group" {
  name        = "Cost Policy Group"
  description = "An approval is necessary once the threshold is crossed"

  control_policies = [
    {
      control_policy_id = data.cm_control_policy.cost_threshold_to_approve.id
      severity          = "high"
    }
  ]
}
