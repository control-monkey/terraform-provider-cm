terraform {
  required_providers {
    cm = {
      source = "example.com/control-monkey/cm"
      version = ">= 1.0"
    }
  }
}

provider "cm" {} # use `export CONTROL_MONKEY_TOKEN=<TOKEN_HERE>` with a valid token



resource "cm_run_task" "firstTask" {
  is_enabled = true
  name       = "AAAA - First Task"
  url        = "https://registry.terraform.io/providers/hashicorp/tfe/latest/docs/resources/organization_run_task"
  hmac_key = "myKey"
}


resource "cm_stack" "TestStack" {
  deployment_behavior = {
    deploy_on_push = true
  }

  namespace_id = "ns-m8ye2gg9k0"
  name         = "TestStack"
  iac_type     = "terraform"

  vcs_info = {
    provider_id = "vcsp-wwvlibsmfw"
    repo_name   = "terraform-test"
    path        = "some/path"
    branch      = "main"
  }

  run_task_config = {
    run_tasks = [
      {
        run_task_id = cm_run_task.firstTask.id
        stage = "preApply"
        enforcement_level = "warning"
      }
    ]
  }
}



resource "cm_stack" "TestStack2" {
  deployment_behavior = {
    deploy_on_push = true
  }

  namespace_id = "ns-m8ye2gg9k0"
  name         = "TestStack2"
  iac_type     = "terraform"

  vcs_info = {
    provider_id = "vcsp-wwvlibsmfw"
    repo_name   = "terraform-test"
    path        = "some/path"
    branch      = "main"
  }

  run_task_config = {
    run_tasks = [
      {
        run_task_id = cm_run_task.firstTask.id
        stage = "postPlan"
        enforcement_level = "warning"
      }
    ]
  }
}