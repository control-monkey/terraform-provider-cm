resource "cm_run_task" "example" {
  name       = "My Run Task"
  url        = "https://my-run-task-endpoint.example.com/callback"
  is_enabled = true
  hmac_key   = "my-secret-hmac-key"
}
