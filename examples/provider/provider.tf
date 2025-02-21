terraform {
    required_providers {
      tines = {
        source = "tines/tines"
        version = "~> 0.2.0"
      }
    }
}

provider "tines" {
    tenant = "https://example.tines.com"
    api_key = var.tines_api_key
}

# Create a Tines Story
resource "tines_story" "example_story" {
    #
}