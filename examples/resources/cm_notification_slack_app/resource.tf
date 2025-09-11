variable "slack_bot_token" {
  type      = string
  sensitive = true
}

resource "cm_notification_slack_app" "slack_app" {
  name           = "example-slack-app"
  bot_auth_token = var.slack_bot_token
}
