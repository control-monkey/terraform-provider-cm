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
