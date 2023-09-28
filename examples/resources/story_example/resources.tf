resource "tines_story" "dev_story_name" {
  data            = file("${path.module}/story-example.json")
  tenant_url      = "https://dev-tenant.tines.com"
  tines_api_token = var.dev_tines_api_token
  team_id         = var.team_id   # optional
  folder_id       = var.folder_id # optional
}


resource "tines_story" "prod_story_name" {
  data            = file("${path.module}/story-example.json")
  tenant_url      = "https://prod-tenant.tines.com"
  tines_api_token = var.prod_tines_api_token
  team_id         = var.team_id   # optional
  folder_id       = var.folder_id # optional
}
