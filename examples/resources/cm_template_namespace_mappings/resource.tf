resource "cm_template_namespace_mappings" "mappings" {
  template_id = cm_template.template.id

  namespaces = [
    {
      namespace_id = cm_namespace.namespace1.id
    },
    {
      namespace_id = cm_namespace.namespace2.id
    }
  ]
}