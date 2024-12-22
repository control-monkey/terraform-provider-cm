resource "cm_disaster_recovery_configuration" "dr_config" {
  scope            = "aws"
  cloud_account_id = "123465789"

  backup_strategy = {
    include_managed_resources = true
    mode                      = "default"

    vcs_info = {
      provider_id = "vcsp-123"
      repo_name   = "terraform-backup"
      branch      = "main"
    }
  }
}