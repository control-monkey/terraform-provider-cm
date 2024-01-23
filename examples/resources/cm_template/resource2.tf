resource "cm_template" "temporary_demo_template" {
  name = "Demo Template"
  iac_type = "terraform"
  description = "Template for temporary demo environment with TTL"

  vcs_info = {
    provider_id = "vcsp-123"
    repo_name = "terraform"
    path = "demo/template"
    branch = "demo"
  }

  policy = {
    ttl_config = {
      max_ttl = {
        type  = "days"
        value = "1"
      }
      default_ttl = {
        type  = "hours"
        value = "5"
      }
    }
  }
}
