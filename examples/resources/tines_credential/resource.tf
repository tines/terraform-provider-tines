resource "tines_credential" "example_text_credential" {
  name    = "Example API Key"
  team_id = 1
  value   = "example-secret-value"
}

# Secret rotation example: wire the value to an upstream secret (e.g. a rotated
# token from your CI/CD pipeline or another Terraform resource). When the value
# changes, Terraform updates the Tines Credential in place during apply.
resource "tines_credential" "rotated_token" {
  name    = "example_token"
  mode    = "TEXT"
  team_id = 2
  value   = var.rotated_git_pat
}

# Write-only example (Terraform >= 1.11): the secret is supplied via `value_wo`
# so it is never written to Terraform state. Because write-only values are not
# tracked for drift, bump `value_wo_version` whenever the secret changes to
# trigger an in-place update (secret rotation).
resource "tines_credential" "write_only_token" {
  name             = "write_only_token"
  team_id          = 2
  value_wo         = var.rotated_git_pat
  value_wo_version = 1
}
