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

  run_trigger = {
    patterns = [
      # matches all files in any sub-directory under the stack's directory.
      # ControlMonkey will replace automatically ${stack_path} with "dev/auto-scaling-group".
      "$${stack_path}/**/*" # '$$' is used to escape the '$' character
    ]
  }

  runner_config = {
    mode = "selfHosted"
    groups = ["dev-runner"]
  }
}