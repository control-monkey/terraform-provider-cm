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
  }

  deployment_approval_policy = {
    rules = [
      {
        type = "requireTwoApprovals"
      }
    ]
  }

  auto_sync = {
    deploy_when_drift_detected = true
  }
}