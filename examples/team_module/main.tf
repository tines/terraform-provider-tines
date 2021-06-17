terraform {
  required_providers {
    tines = {
      source  = "github.com/tuckner/tines"
      version = "0.0.18"
    }
  }
}

provider "tines" {
  alias = "t"
}

module "team_environment" {
  providers = {
    tines = tines.t
  }
  for_each = var.users
  source = "./teams"
  user = each.value.name
  email = each.value.email
}