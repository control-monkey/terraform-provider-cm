resource "cm_events_subscriptions" "events_subscriptions" {
  scope         = "namespace"
  scope_id      = cm_namespace.namespace.id
  subscriptions = [
    {
      event_type               = "stack::deployment::applyStarted"
      notification_endpoint_id = cm_notification_endpoint.endpoint1.id
    },
    {
      event_type               = "stack::deployment::applyFinished"
      notification_endpoint_id = cm_notification_endpoint.endpoint1.id
    },
    {
      event_type               = "stack::deployment::failed"
      notification_endpoint_id = cm_notification_endpoint.endpoint2.id
    }
  ]
}

