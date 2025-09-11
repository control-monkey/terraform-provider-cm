resource "cm_notification_endpoint" "notification_endpoint" {
  name = "ControlMonkey Notifications"
  protocol = "slack"
  url = "https://hooks.slack.com"
}
