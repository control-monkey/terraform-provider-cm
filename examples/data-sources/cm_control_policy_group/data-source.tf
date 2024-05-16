data "cm_control_policy_group" "PCI_DSS" {
  name = "PCI DSS"
}


resource "cm_control_policy_group_mappings" "PCI_DSS_policy_group_mappings" {
  control_policy_group_id = data.cm_control_policy_group.PCI_DSS.id

  targets = [
    {
      target_id         = "ALL"
      target_type       = "namespace"
      enforcement_level = "hardMandatory"
    }
  ]
}
