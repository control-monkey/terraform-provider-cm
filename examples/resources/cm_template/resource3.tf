resource "cm_template" "rds_template" {
  name     = "Ephemeral RDS For R&D"
  iac_type = "terraform"

  vcs_info = {
    provider_id = "vcsp-123"
    repo_name   = "terraform"
    path        = "templates/rds"
  }

  policy = {
    ttl_config = {
      max_ttl = {
        type  = "days"
        value = "10"
      }
      default_ttl = {
        type  = "days"
        value = "5"
      }
    }
  }
}

resource "cm_variable" "creator_variable" {
  scope          = "template"
  scopeId        = cm_template.rds_template.id
  key            = "creator_name"
  type           = "tfVar"
  is_sensitive   = false
  is_overridable = true
  is_required    = true
}

resource "cm_variable" "allowed_instances_variable" {
  scope          = "template"
  scopeId        = cm_template.rds_template.id
  key            = "instance_type"
  type           = "tfVar"
  is_sensitive   = false
  is_overridable = true
  is_required    = true
}