terraform {
  required_providers {
    tines = {
      source  = "github.com/tuckner/tines"
      version = "0.0.18"
    }
  }
}

provider "tines" {}

provider "aws" {
  region    = "us-west-1"
}

## Create Tines Resources

resource "tines_credential" "aws_key" {
    name = "aws_key"
    mode = "AWS"
    aws_authentication_type = "KEY"
    aws_access_key = "replace"
    aws_secret_key = "me"
    team_id = var.tines_team
}

resource "random_id" "webhook_secret" {
  byte_length = 8
}

## Create Tines Story & Actions

resource "tines_story" "aws_response" {
    name = "AWS Response"
    team_id = var.tines_team
}

resource "tines_agent" "webhook" {
  name = "Receive SNS"
  agent_type = "Agents::WebhookAgent"
  story_id = tines_story.aws_response.id
  keep_events_for = 604800
  source_ids = []
  receiver_ids = [tines_agent.if_subscription_registration.id, tines_agent.expand_event_3.id]
  agent_options = jsonencode({
    "secret": random_id.webhook_secret.dec,
    "verbs": "get,post"
  })
  position = {
    x = 0.0
    y = 0.0
  }
}

resource "tines_agent" "if_subscription_registration" {
    name = "If Subscription Registration"
    agent_type = "Agents::TriggerAgent"
    story_id = tines_story.aws_response.id
    keep_events_for = 0
    source_ids = []
    receiver_ids = [tines_agent.webhook_confirm.id]
    position = {
      x = 120.0
      y = 75.0
    }
    agent_options = jsonencode({"rules": [{"path": "{{ .receive_sns.body.SubscribeURL }}", "type": "!regex", "value": "^$"}]})
}

resource "tines_agent" "webhook_confirm" {
  name = "Webhook Confirm"
  agent_type = "Agents::HTTPRequestAgent"
  story_id = tines_story.aws_response.id
  keep_events_for = 604800
  source_ids = []
  receiver_ids = []
  agent_options = jsonencode({
    "url": "{{ .receive_sns.body.SubscribeURL }}",
    "method": "get"
  })
  position = {
    x = 120.0
    y = 165.0
  }
}

resource "tines_agent" "expand_event_3" {
    name = "Expand Event"
    agent_type = "Agents::EventTransformationAgent"
    story_id = tines_story.aws_response.id
    keep_events_for = 0
    source_ids = []
    receiver_ids = [tines_agent.is_bucketanonymousaccessgranted_4.id]
    position = {
      x = -120.0
      y = 75.0
    }
    agent_options = jsonencode({"mode": "message_only", "payload": {"message": "{{ .receive_sns.body.Message | as_object }}"}})
}

resource "tines_agent" "is_bucketanonymousaccessgranted_4" {
    name = "Is BucketAnonymousAccessGranted"
    agent_type = "Agents::TriggerAgent"
    story_id = tines_story.aws_response.id
    keep_events_for = 0
    source_ids = []
    receiver_ids = [tines_agent.apply_aws_s3_bucket_block_policy_5.id]
    position = {
      x = -120.0
      y = 150.0
    }
    agent_options = jsonencode({"rules": [{"path": "{{ .expand_event.message.detail.type }}", "type": "field==value", "value": "Policy:S3/BucketAnonymousAccessGranted"}]})
}

resource "tines_agent" "apply_aws_s3_bucket_block_policy_5" {
    name = "Apply AWS S3 Bucket Block Policy"
    agent_type = "Agents::HTTPRequestAgent"
    story_id = tines_story.aws_response.id
    keep_events_for = 0
    source_ids = []
    receiver_ids = []
    position = {
      x = -120.0
      y = 255.0
    }
    agent_options = jsonencode({"content_type": "application/xml", "headers": {"Authorization": "{{ .CREDENTIAL.aws_key }}"}, "method": "put", "payload": "\u003c?xml version=\"1.0\" encoding=\"UTF-8\"?\u003e\n\u003cPublicAccessBlockConfiguration\u003e\n      \u003cBlockPublicAcls\u003eTRUE\u003c/BlockPublicAcls\u003e\n      \u003cIgnorePublicAcls\u003eTRUE\u003c/IgnorePublicAcls\u003e\n      \u003cBlockPublicPolicy\u003eTRUE\u003c/BlockPublicPolicy\u003e\n      \u003cRestrictPublicBuckets\u003eTRUE\u003c/RestrictPublicBuckets\u003e\n\u003c/PublicAccessBlockConfiguration\u003e", "url": "https://{{.expand_event.message.detail.resource.s3BucketDetails.first.name}}.s3.{{.expand_event.message.detail.region}}.amazonaws.com/?publicAccessBlock"})
    depends_on = [
      tines_credential.aws_key,
    ]
}

## Setup Guardduty and CloudWatch

resource "aws_guardduty_detector" "guardduty" {
  enable = true
  finding_publishing_frequency = "FIFTEEN_MINUTES"
}

resource "aws_cloudwatch_event_rule" "main" {
  name          = "guardduty-finding-events"
  description   = "AWS GuardDuty event findings"
  event_pattern = <<EOF
{
  "source": [
    "aws.guardduty"
  ]
}
EOF
}

## Setup SNS

resource "aws_sns_topic" "guardduty_sns" {
  name = "guardduty-sns"
}

resource "aws_cloudwatch_event_target" "sns" {
  rule      = aws_cloudwatch_event_rule.main.name
  target_id = "send-to-guardduty-sns"
  arn       = aws_sns_topic.guardduty_sns.arn
}

resource "aws_sns_topic_subscription" "sns_target" {
  depends_on = [tines_agent.webhook_confirm]
  topic_arn = aws_sns_topic.guardduty_sns.arn
  protocol = "https"
  endpoint = format("%s/webhook/%s/%s", var.tines_base_url, tines_agent.webhook.guid, random_id.webhook_secret.dec)
  endpoint_auto_confirms = true
}
