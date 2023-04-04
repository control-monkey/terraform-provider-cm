terraform {
  required_providers {
    cm = {
      source  = "controlmonkey/cm"
      version = ">= 1.0"
    }
  }
}

provider "cm" {
  // You can also set this via CONTROL_MONKEY_TOKEN environment variable.
  token = "CONTROL_MONKEY_TOKEN"
}
