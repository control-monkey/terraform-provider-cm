data "cm_blueprint" "demo_env" {
  name = "Demo Env Blueprint"
}

data "cm_namespace" "dev_namespace" {
  name = "Dev"
}


resource "cm_blueprint_namespace_mappings" "demo_blueprint_namespace_mappings" {
  blueprint_id = data.cm_blueprint.demo_env.id

  namespaces = [
    {
      namespace_id = data.cm_namespace.dev_namespace.id
    }
  ]
}
