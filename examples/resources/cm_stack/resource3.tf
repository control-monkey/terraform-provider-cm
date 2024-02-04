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

  iac_config = {
    terraform_version = "1.4.5"
  }

  policy = {
    ttl_config = {
      ttl = {
        type  = "hours"
        value = "3"
      }
    }
  }
}