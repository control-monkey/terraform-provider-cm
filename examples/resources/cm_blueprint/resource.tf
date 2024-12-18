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
      provider_id  = "vcsp-github"
      repo_name    = "terraform"
      path_pattern = "infra/{env}/{region}/{service}"
    }
  }

  substitute_parameters = [
    {
      key         = "env"
      description = "The name of the environment e.g prod, stage"
    },
    {
      key         = "region"
      description = "The region in the cloud e.g us-east-1"
    },
    {
      key         = "service"
      description = "The name of the service e.g ec2, s3"
    }
  ]
}