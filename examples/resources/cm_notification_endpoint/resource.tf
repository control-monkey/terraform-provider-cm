resource "cm_notification_endpoint" "notification_endpoint" {
  name = "ControlMonkey Notifications"
  protocol = "slack"
  url = "https://www.slack.com/example/webhook"
}
