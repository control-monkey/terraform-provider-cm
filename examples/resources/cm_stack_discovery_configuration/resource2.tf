resource "cm_namespace" "dev_namespace" {
  name        = "Development"
  description = "Development environment namespace"
}

resource "cm_stack_discovery_configuration" "dev_discovery" {
  name         = "Development Stack Discovery"
  namespace_id = cm_namespace.dev_namespace.id
  description  = "Auto-discover development stacks with self-hosted runners"

  vcs_patterns = [
    {
      provider_id = "vcsp-def456"
      repo_name   = "my-org/dev-infrastructure"
      path_patterns = ["apps/*/terraform/**", "services/*/infra/**"]
      exclude_path_patterns = ["**/node_modules/**", "**/vendor/**"]
      branch      = "develop"
    },
    {
      provider_id = "vcsp-def456"
      repo_name   = "my-org/microservices"
      path_patterns = ["**/infrastructure/**"]
      branch      = "main"
    }
  ]

  stack_config = {
    iac_type = "terraform"

    deployment_behavior = {
      deploy_on_push = true
    }

    iac_config = {
      terraform_version = "1.5.5"
    }

    runner_config = {
      mode = "selfHosted"
      groups = ["dev-runners", "shared-runners"]
    }

    auto_sync = {
      deploy_when_drift_detected = false
    }
  }
}
