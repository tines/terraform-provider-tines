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

# resource "tines-v2_story" "prod_story_name" {
#   data            = file("${path.module}/prod-story.json")
#   file_hash       = filesha512("${path.module}/prod-story.json")
#   tenant_url      = var.prod_tenant_url
#   tines_api_token = var.prod_tines_api_token
# }
