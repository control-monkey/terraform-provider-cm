data "cm_namespace" "namespace" {
  name = "Production"
}

resource "cm_stack_discovery_configuration" "example" {
  name         = "Terraform Stack Discovery"
  namespace_id = data.cm_namespace.namespace.id
  description  = "Auto-discover and manage Terraform stacks in production repositories"

  vcs_patterns = [
    {
      provider_id = "vcsp-abc123"
      repo_name   = "my-org/infrastructure"
      path_patterns = ["environments/*/terraform/**", "modules/*/**"]
      exclude_path_patterns = ["**/test/**", "**/.terraform/**"]
      branch      = "main"
    }
  ]

  stack_config = {
    iac_type = "terraform"

    deployment_behavior = {
      deploy_on_push = false
    }

    deployment_approval_policy = {
      rules = [
        {
          type = "requireTwoApprovals"
        }
      ]
    }

    iac_config = {
      terraform_version = "1.5.5"
    }

    runner_config = {
      mode = "managed"
    }

    auto_sync = {
      deploy_when_drift_detected = true
    }
  }
}

