resource "cm_blueprint" "blueprint" {
  name = "Blueprint example"

  blueprint_vcs_info = {
    provider_id = "vcsp-github"
    repo_name   = "terraform"
    path        = "infra/blueprints/vpc"
  }

  stack_configuration = {
    name_pattern = "{env}-{region}-{service}"
    iac_type     = "terraform"

    vcs_info_with_patterns = {
      provider_id    = "vcsp-github"
      repo_name      = "terraform-developers"
      path_pattern   = "infra/{env}/{region}/{service}"
      branch_pattern = "{env}-development"
    }

    deployment_approval_policy = {
      rules = [
        {
          type = "requireTwoApprovals"
        }
      ]
    }
  }

  substitute_parameters = [
    {
      key         = "env"
      description = "The name of the environment in lower case e.g prod, stage"
      value_conditions = [
        {
          operator = "in"
          values = ["prod", "stage", "dev"]
        }
      ]
    },
    {
      key         = "region"
      description = "The region in the cloud e.g us-east-1"
    },
    {
      key         = "service"
      description = "The name of the service in lower case e.g ec2, s3"
    }
  ]
}