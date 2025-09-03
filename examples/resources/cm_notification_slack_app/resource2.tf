resource "cm_notification_slack_app" "slack_app" {
  name           = "my-slack-app"
  bot_auth_token = "ignored-token"

  lifecycle {
    ignore_changes = [bot_auth_token]
  }
}
