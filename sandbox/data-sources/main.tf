terraform {
  required_providers {
    cm = {
      source = "example.com/control-monkey/cm"
      version = ">= 1.0"
    }
  }
}

provider "cm" {} # use `export CONTROL_MONKEY_TOKEN=<TOKEN_HERE>` with a valid token
