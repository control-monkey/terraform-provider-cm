---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cm_notification_endpoint Data Source - terraform-provider-cm"
subcategory: ""
description: |-
  
---

# cm_notification_endpoint (Data Source)



## Example Usage

```terraform
data "cm_namespace" "prod_namespace" {
  name = "Prod"
}

data "cm_notification_endpoint" "controlmonkey_notification_endpoint" {
  name = "ControlMonkey Notifications"
}


resource "cm_events_subscriptions" "prod_events_subscriptions" {
  scope         = "namespace"
  scope_id      = data.cm_namespace.prod_namespace.id
  subscriptions = [
    {
      event_type               = "stack::deployment::applyStarted"
      notification_endpoint_id = data.cm_notification_endpoint.controlmonkey_notification_endpoint.id
    },
    {
      event_type               = "stack::deployment::applyFinished"
      notification_endpoint_id = data.cm_notification_endpoint.controlmonkey_notification_endpoint.id
    },
    {
      event_type               = "stack::deployment::failed"
      notification_endpoint_id = data.cm_notification_endpoint.controlmonkey_notification_endpoint.id
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) The unique ID of the notification endpoint.
- `name` (String) The name of the notification endpoint.
