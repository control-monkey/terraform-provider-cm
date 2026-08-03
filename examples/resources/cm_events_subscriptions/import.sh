# The ID is "<scope>/<scope_id>". Only the organization scope has no scope_id.
terraform import cm_events_subscriptions.events_subscriptions namespace/ns-123
terraform import cm_events_subscriptions.events_subscriptions namespace/ALL
terraform import cm_events_subscriptions.events_subscriptions stack/stk-123
terraform import cm_events_subscriptions.events_subscriptions awsAccount/123456789012
terraform import cm_events_subscriptions.events_subscriptions azureSubscription/00000000-1111-2222-3333-444444444444
terraform import cm_events_subscriptions.events_subscriptions gcpProject/my-gcp-project
terraform import cm_events_subscriptions.events_subscriptions organization
