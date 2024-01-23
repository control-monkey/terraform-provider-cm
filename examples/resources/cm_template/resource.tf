resource "cm_template" "template_developers" {
  name = "Dev Self-Service Template"
  iac_type = "terraform"
  description = "Self service on Dev environment for developers"

  vcs_info = {
    provider_id = "vcsp-123"
    repo_name = "terraform"
    path = "dev/self-service"
  }
}