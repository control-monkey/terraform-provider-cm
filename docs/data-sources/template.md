---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cm_template Data Source - terraform-provider-cm"
subcategory: ""
description: |-
  
---

# cm_template (Data Source)



## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) The unique ID of the template.
- `name` (String) The name of the template.