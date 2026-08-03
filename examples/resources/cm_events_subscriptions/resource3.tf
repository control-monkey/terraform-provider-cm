# Organization scope covers the whole organization and takes NO `scope_id`. Only the three
# console-operation event types are accepted here.
resource "cm_events_subscriptions" "organization" {
  scope = "organization"
  subscriptions = [
    {
      event_type               = "aws::consoleOperation"
      notification_endpoint_id = cm_notification_endpoint.endpoint1.id
    },
    {
      event_type               = "azure::consoleOperation"
      notification_endpoint_id = cm_notification_endpoint.endpoint1.id
    },
    {
      event_type               = "gcp::consoleOperation"
      notification_endpoint_id = cm_notification_endpoint.endpoint2.id
    }
  ]
}

# GCP project scope. `scope_id` is the GCP project ID, not a ControlMonkey ID. Only
# `gcp::consoleOperation` is accepted.
resource "cm_events_subscriptions" "gcp_project" {
  scope    = "gcpProject"
  scope_id = "my-gcp-project-123"
  subscriptions = [
    {
      event_type               = "gcp::consoleOperation"
      notification_endpoint_id = cm_notification_endpoint.endpoint1.id
    }
  ]
}

# Stack scope. Accepts the stack, plan and runner event types - not the console-operation ones.
resource "cm_events_subscriptions" "stack" {
  scope    = "stack"
  scope_id = cm_stack.stack.id
  subscriptions = [
    {
      event_type               = "stack::driftDetected"
      notification_endpoint_id = cm_notification_endpoint.endpoint1.id
    },
    {
      event_type               = "stack::plan::failed"
      notification_endpoint_id = cm_notification_endpoint.endpoint2.id
    }
  ]
}
