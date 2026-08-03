# Keep the flow definition in your repository and load it from a file
resource "cm_custom_flow" "from_file" {
  name       = "Custom flow from file"
  is_enabled = true
  flow_yaml  = file("${path.module}/cm.yaml")
}
