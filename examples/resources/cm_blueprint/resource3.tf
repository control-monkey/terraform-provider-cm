data "cm_namespace" "prod_namespace" {
  name = "Prod"
}

data "cm_team" "prod_devops_team" {
  name = "Prod DevOps Team"
}

resource "cm_blueprint_namespace_mappings" "vpc_blueprint_namespace_mappings" {
  blueprint_id = cm_blueprint.blueprint.id

  namespaces = [
    {
      namespace_id = data.cm_namespace.prod_namespace.id
    }
  ]
}

resource "cm_blueprint" "blueprint" {
  name = "Blueprint example"

  blueprint_vcs_info = {
    provider_id = "vcsp-github"
    repo_name   = "terraform"
    path        = "infra/blueprints/vpc"
    branch      = "cm-blueprints"
  }

  stack_configuration = {
    name_pattern = "{env}-{region}-{service}"
    iac_type     = "terraform"

    vcs_info_with_patterns = {
      provider_id    = "vcsp-github"
      repo_name      = "terraform"
      path_pattern   = "infra/{env}/{region}/{service}"
      branch_pattern = "{env}-main"
    }

    deployment_approval_policy = {
      rules = [
        {
          type = "requireTeamsApproval"
          parameters = jsonencode({
            teams = [data.cm_team.prod_devops_team.id]
          })
        },
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
    },
    {
      key         = "code_dynamic_parameter"
      description = "The provided value will replace all occurrences of {code_dynamic_parameter} in the code"
    }
  ]
}
