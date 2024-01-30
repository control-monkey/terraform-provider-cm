resource "cm_control_policy_mapping" "no_public_bucket_prod" {
  control_policy_id = "pol-123"
  target_id         = cm_namespace.prod_namespace.id
  target_type       = "namespace"
  enforcement_level = "hardMandatory"
}