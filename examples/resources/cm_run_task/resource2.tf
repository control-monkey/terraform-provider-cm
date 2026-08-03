variable "run_task_hmac_key" {
  type      = string
  sensitive = true
}

resource "cm_run_task" "example" {
  name       = "My Run Task"
  url        = "https://my-run-task-endpoint.example.com/callback"
  is_enabled = true

  hmac_key_wo         = var.run_task_hmac_key
  hmac_key_wo_version = 1
}
