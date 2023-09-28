terraform {
  required_providers {
    tines = {
      source = "tines/tines"
    }
  }
}

provider "tines" {}

resource "tines_story" "dev_story_name" {
  data            = file("${path.module}/dev-story.json")
  tenant_url      = var.dev_tenant_url
  tines_api_token = var.dev_tines_api_token
  team_id         = var.team_id # optional
  # folder_id       = var.folder_id # optional
}
