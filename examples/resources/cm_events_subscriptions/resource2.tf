# Cloud account scopes. `scope_id` is the cloud provider's own identifier, not a ControlMonkey ID,
# and each cloud account can only subscribe to its own console-operation event type.

resource "cm_events_subscriptions" "aws_account" {
  scope    = "awsAccount"
  scope_id = "123456789012" # the 12 digit AWS account ID
  subscriptions = [
    {
      event_type               = "aws::consoleOperation"
      notification_endpoint_id = cm_notification_endpoint.endpoint1.id
    }
  ]
}

resource "cm_events_subscriptions" "azure_subscription" {
  scope    = "azureSubscription"
  scope_id = "00000000-1111-2222-3333-444444444444" # the Azure subscription ID
  subscriptions = [
    {
      event_type               = "azure::consoleOperation"
      notification_endpoint_id = cm_notification_endpoint.endpoint1.id
    }
  ]
}

# Every namespace in the organization. A specific-namespace resource does not absorb these.
resource "cm_events_subscriptions" "all_namespaces" {
  scope    = "namespace"
  scope_id = "ALL"
  subscriptions = [
    {
      event_type               = "stack::deployment::failed"
      notification_endpoint_id = cm_notification_endpoint.endpoint2.id
    }
  ]
}
