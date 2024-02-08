resource "cm_blueprint_namespace_mappings" "mappings" {
  blueprint_id = "blp-123"

  namespaces = [
    {
      namespace_id = cm_namespace.namespace1.id
    },
    {
      namespace_id = cm_namespace.namespace2.id
    }
  ]
}