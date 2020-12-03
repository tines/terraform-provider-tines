terraform {
  required_providers {
    tines = {
      source = "github.com/tuckner/tines"
      version = ">=0.0.5"
    }
  }
}

provider "tines" {
  email    = var.tines_email
  base_url = var.tines_base_url
  token    = var.tines_token
}

resource "random_id" "webhook_secret" {
  byte_length = 8
}

resource "tines_global_resource" "ip_list" {
  name = "ip_list"
  value_type = "text"
  value = "10.0.0.0\n10.0.0.1\n"
}

resource "tines_agent" "webhook" {
  name = "Webhook Agent"
  agent_type = "Agents::WebhookAgent"
  story_id = var.story_id
  keep_events_for = 604800
  source_ids = []
  receiver_ids = []
  agent_options = jsonencode({
    "secret": random_id.webhook_secret.dec,
    "verbs": "get,post",
    "response": format("{{.RESOURCE.%s}}", tines_global_resource.ip_list.name)
  })
}

