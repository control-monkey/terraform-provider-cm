---
page_title: "cm_events_subscriptions Resource - terraform-provider-cm"
subcategory: ""
description: |-
  Creates, updates and destroys events subscriptions.
---

# cm_events_subscriptions (Resource)

Creates, updates and destroys events subscriptions.

## Learn More

- [Email Alerts for IaC Events in ControlMonkey](https://controlmonkey.io/news/controlmonkey-email-alerts/)
- [Console Operations Notifications](https://controlmonkey.io/news/console-operations-notifications/)
- [Terraform Microsoft Teams Support: Real-Time Infrastructure Notifications](https://controlmonkey.io/news/teams-notification-support/)

## Example Usage

### Namespace scope
```terraform
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
```

### Cloud account scopes and all namespaces
```terraform
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
```

### Organization, GCP project and stack scopes
```terraform
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
```

Each scope accepts only certain `event_type` values, and the API rejects the rest:

| `scope` | `scope_id` | Allowed `event_type` |
| --- | --- | --- |
| `organization` | must not be set | `aws::consoleOperation`, `azure::consoleOperation`, `gcp::consoleOperation` |
| `awsAccount` | the 12 digit AWS account ID | `aws::consoleOperation` |
| `azureSubscription` | the Azure subscription ID | `azure::consoleOperation` |
| `gcpProject` | the GCP project ID | `gcp::consoleOperation` |
| `stack` | a ControlMonkey stack ID | the `stack::*`, `stack::plan::*` and `runner::*` types |
| `namespace` | a namespace ID, or `ALL` | the stack types, plus `stack::created` and `stack::createdByAutoDiscovery` |

Each resource owns only the subscriptions attached to its exact `scope` and `scope_id`. A
`scope_id` of `ALL`, or an `organization` scoped subscription, is inherited by narrower scopes at
runtime but is not managed by their resources - so a specific namespace resource and an `ALL`
resource can coexist without fighting over the same rows.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `scope` (String) Scope of the resource. Known values: [organization, namespace, awsAccount, azureSubscription, gcpProject, stack].

### Optional

- `scope_id` (String) The ID of the resource to which the subscriptions are attached. Must not be set when `scope` is `organization`, and is required for every other scope. For `namespace` use the namespace ID, or `ALL` to target every namespace in the organization. For `stack` use the stack ID. For `awsAccount`, `azureSubscription` and `gcpProject` use the cloud provider's own identifier - the 12 digit AWS account ID, the Azure subscription ID or the GCP project ID - not a ControlMonkey ID. ControlMonkey does not validate that a cloud identifier exists, so a typo creates a subscription that never fires.
- `subscriptions` (Attributes Set) Specifies a list of events subscriptions. (see [below for nested schema](#nestedatt--subscriptions))

### Read-Only

- `id` (String) The unique ID of this resource.

<a id="nestedatt--subscriptions"></a>
### Nested Schema for `subscriptions`

Required:

- `event_type` (String) The type of the event. Not every event type is valid for every `scope` - cloud console operations belong to the cloud account and organization scopes, while stack and plan events belong to the stack and namespace scopes. ControlMonkey rejects an invalid combination and lists the allowed types for the scope in the error. Find supported types [here](https://docs.controlmonkey.io/controlmonkey-api/api-enumerations#event-types)
- `notification_endpoint_id` (String) The unique ID of the endpoint to which the notification will be sent.

Read-Only:

- `id` (String) The unique ID of the subscription.

## Import

`cm_events_subscriptions` can be imported using the following format `scope/scope_id` or only `scope` if scope_id does not exist, e.g.

```shell
# The ID is "<scope>/<scope_id>". Only the organization scope has no scope_id.
terraform import cm_events_subscriptions.events_subscriptions namespace/ns-123
terraform import cm_events_subscriptions.events_subscriptions namespace/ALL
terraform import cm_events_subscriptions.events_subscriptions stack/stk-123
terraform import cm_events_subscriptions.events_subscriptions awsAccount/123456789012
terraform import cm_events_subscriptions.events_subscriptions azureSubscription/00000000-1111-2222-3333-444444444444
terraform import cm_events_subscriptions.events_subscriptions gcpProject/my-gcp-project
terraform import cm_events_subscriptions.events_subscriptions organization
```