# Manage this Story resource using a JSON story export.
resource "tines_story" "example_imported_story" {
  team_id = 1
  data = file("${path.module}/story-example.json")
}

# Manage this Story resource using resource attributes.
resource "tines_story" "example_configured_story" {
  team_id = 1
  name = "Example Story Name"
  change_control_enabled = true
}