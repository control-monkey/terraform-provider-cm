resource "cm_stack" "auto_scaling_group_dev" {
  name         = "dev/auto-scaling-group"
  description  = "Auto Scaling Group Stack"
  namespace_id = "ns-dev"
  iac_type     = "terraform"

  deployment_behavior = {
    deploy_on_push = true
  }

  vcs_info = {
    provider_id = "vcsp-github"
    repo_name   = "terraform"
    path        = "dev/auto-scaling-group"
  }

  run_trigger = {
    patterns = [
      "${stack_path}/**/*" # matches all files in any sub-directory under the stack's directory.
    ] # ControlMonkey will replace automatically ${stack_path} with "dev/auto-scaling-group"
  }

  runner_config = {
    mode = "selfHosted"
    groups = ["dev-runner"]
  }
}