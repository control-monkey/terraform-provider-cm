resource "cm_run_task_mappings" "security_scan_mappings" {
  run_task_id = cm_run_task.security_scan.id

  targets = [
    {
      target_id         = cm_namespace.namespace.id
      target_type       = "namespace"
      enforcement_level = "hardMandatory"
      stage             = "postPlan"
    },
    {
      target_id         = cm_stack.stack.id
      target_type       = "stack"
      enforcement_level = "softMandatory"
      stage             = "postPlan"
    },
  ]
}
