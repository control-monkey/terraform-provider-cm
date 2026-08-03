# Map a run task to every namespace in the organization
resource "cm_run_task_mappings" "all_namespaces" {
  run_task_id = cm_run_task.security_scan.id

  targets = [
    {
      target_id         = "ALL"
      target_type       = "namespace"
      enforcement_level = "warning"
      stage             = "postPlan"
    },
  ]
}
