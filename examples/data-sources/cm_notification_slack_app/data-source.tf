data "cm_notification_slack_app" "app" {
  name = "Slack"
}

resource "cm_notification_endpoint" "slack_app" {
  name     = "example-slack-app-endpoint"
  protocol = "slackApp"

  slack_app_config = {
    notification_slack_app_id = data.cm_notification_slack_app.app.id
    channel_id                = "C0123456789"
  }
}
