terraform {
  required_providers {
    tines = {
      source = "github.com/tuckner/tines"
      version = ">= 0.0.5"
    }
  }
}

provider "tines" {
  email    = var.tines_email
  base_url = var.tines_base_url
  token    = var.tines_token
}

provider "aws" {
  region    = "us-east-1"
}

resource "random_id" "webhook_secret" {
  byte_length = 8
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
    "verbs": "get,post"
  })
}

resource "tines_agent" "webhook_confirm" {
  name = "Webhook Confirm"
  agent_type = "Agents::HTTPRequestAgent"
  story_id = var.story_id
  keep_events_for = 604800
  source_ids = [tines_agent.webhook.agent_id]
  receiver_ids = []
  agent_options = jsonencode({
    "url": "{{ .webhook_agent.body.SubscribeURL }}",
    "method": "get"
  })
}

resource "aws_sns_topic_subscription" "sns_target" {
  depends_on = [tines_agent.webhook_confirm]
  topic_arn = var.aws_topic_arn
  protocol = "https"
  endpoint = format("%s/webhook/%s/%s", var.tines_base_url, tines_agent.webhook.guid, random_id.webhook_secret.dec)
  endpoint_auto_confirms = true
}
