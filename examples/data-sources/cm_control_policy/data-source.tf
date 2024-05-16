data "cm_namespace" "prod_namespace" {
  name = "Prod"
}

data "cm_control_policy" "allowed_aws_regions" {
  name = "Allowed AWS Regions"
}


resource "cm_control_policy_mappings" "allowed_regions_policy_mappings" {
  control_policy_id = data.cm_control_policy.allowed_aws_regions.id

  targets = [
    {
      target_id         = data.cm_namespace.prod_namespace.id
      target_type       = "namespace"
      enforcement_level = "hardMandatory"
    }
  ]
}
