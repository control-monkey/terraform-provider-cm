resource "cm_variable" "ephemeral_stack_instance_types" {
  scope           = "template"
  scope_id        = "tmpl-instances"
  key             = "instance_type"
  type            = "tfVar"
  description     = "This tfVar is injected as the instance type of the EC2 instance"
  is_sensitive    = false
  is_overridable  = true
  is_required     = true
  value_conditions = [
    {
      operator = "in"
      values   = ["t2.micro", "t2.nano"]
    }
  ]
}