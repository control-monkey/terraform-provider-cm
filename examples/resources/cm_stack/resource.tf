resource "cm_stack" "stack" {
  iac_type     = "terraform"
  namespace_id = "ns-123"
  name         = "stack1"
  description  = "first stack test"

  deployment_behavior = {
    deploy_on_push    = "true"
    wait_for_approval = "true"
  }

  vcs_info = {
    provider_id = "vcsp-123"
    repo_name   = "nonExistRepo"
  }

  iac_config = {
    terraform_version  = "1.4.5"
    terragrunt_version = "0.45.3"
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