resource "cm_control_policy" "control_policy" {
  name        = "AWS Resources should have the Env tag with value Dev/Stage/Prod"
  description = "All AWS infrastructure should have the Env tag with value Dev/Stage/Prod."
  type        = "aws_required_tags"
  parameters  = jsonencode({
    tags = [
      {
        key           = "Env"
        allowedValues = [
          "Dev",
          "Stage",
          "Prod"
        ]
      }
    ]
  })
}