terraform {
  required_providers {
    tines = {
      source  = "github.com/tuckner/tines"
      version = "0.0.18"
    }
  }
}

resource "tines_team" "team" {
    name = "${var.user}'s Team"
}

resource "null_resource" "invite" {
  provisioner "local-exec" {
    command = "curl -XPOST \"$TINES_URL/api/v1/teams/${tines_team.team.id}/invite_member\" -H \"x-user-email: $TINES_EMAIL\" -H \"x-user-token: $TINES_TOKEN\" -d '{\"email\": \"${var.email}\"}' -H 'content-type: application/json'"
  }
}

