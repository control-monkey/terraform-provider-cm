data "cm_stack" "stack" {
  name = "Stack"
}

data "cm_stack" "depends_on_stack" {
  name = "Depends on Stack"
}

resource "cm_stack_dependency" "dependency" {
  stack_id             = data.cm_stack.stack.id
  depends_on_stack_id  = data.cm_stack.depends_on_stack.id
  trigger_option       = "always"

  references = [
    {
      output_of_stack_to_depend_on = "db_endpoint"
      input_for_stack              = "db_endpoint"
      include_sensitive_output     = false
    }
  ]
}
