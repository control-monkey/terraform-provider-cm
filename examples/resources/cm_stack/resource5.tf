resource "cm_stack" "auto_scaling_group_dev" {
  name         = "dev/auto-scaling-group"
  description  = "Auto Scaling Group Stack"
  namespace_id = cm_namespace.dev_namespace.id
  iac_type     = "terraform"

  deployment_behavior = {
    deploy_on_push = true
  }

  vcs_info = {
    provider_id = "vcsp-123"
    repo_name   = "terraform"
    path        = "dev/auto-scaling-group"
  }

  run_task_config = {
    run_tasks = [
      {
        run_task_id = "rtsk-123"
        stage = "postPlan"
        enforcement_level = "warning"
      }
    ]
  }
}