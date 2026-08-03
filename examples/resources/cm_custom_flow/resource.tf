resource "cm_custom_flow" "opa_validation" {
  name       = "OPA validation"
  is_enabled = true

  flow_yaml = <<-EOT
    version: 1
    customFlows:
      run:
        steps:
          terraformPlan:
            after:
              - name: evaluate OPA policy
                cmd: |
                  opa eval --data policy.rego --input showPlanOutput.json "data.terraform.deny"
                failureBehavior: stop
  EOT
}
