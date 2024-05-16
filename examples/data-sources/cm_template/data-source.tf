data "cm_template" "temporary_workload" {
  name = "Temporary Workload"
}

data "cm_namespace" "dev_namespace" {
  name = "Dev"
}


resource "cm_template_namespace_mappings" "temporary_workload_namespace_mappings" {
  template_id = data.cm_template.temporary_workload.id

  namespaces = [
    {
      namespace_id = data.cm_namespace.dev_namespace.id
    }
  ]
}
