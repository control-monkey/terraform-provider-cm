data "cm_namespace" "prod_namespace" {
  name = "Prod"
}

data "cm_stack" "prod_eks" {
  name         = "EKS"
  namespace_id = data.cm_namespace.prod_namespace.id
}

data "cm_control_policy" "allowed_eks_regions" {
  name = "Allowed EKS Regions"
}


resource "cm_control_policy_mappings" "allowed_eks_regions_policy_mappings" {
  control_policy_id = data.cm_control_policy.allowed_eks_regions.id

  targets = [
    {
      target_id         = data.cm_stack.prod_eks.id
      target_type       = "stack"
      enforcement_level = "hardMandatory"
    }
  ]
}
