---
layout: ""
page_title: "Upgrading to version 0.1.x (from 0.0.x)"
description: Terraform Tines Provider Version 0.1 Upgrade Guide
---

# Terraform Tines Provider Version 0.1 Upgrade Guide
Starting with version `0.1.0`, the Tines provider on Terraform introduces a new and simplified way to manage stories. This means that any Resource and its Schema with a version below `0.1.0` will no longer be compatible with `tines` Terraform provider version `0.1.0` or higher, as there are breaking changes in this version. If you have an older Terraform state file for a version below `0.1.0`, we recommend starting fresh by initializing a new state (via `terraform init`) for version `0.1.0` or higher.


## Provider Version Configuration
If you are not ready to make a move to version 0.1 of the Tines provider, you may keep the 0.0.x branch active for
your Terraform project by specifying:

```terraform
provider "tines" {
  version = "~> 0.0"
  # ... any other configuration
}
```

## Getting Started With v0.1.0
To export your Tines Story, follow [these instructions](https://www.tines.com/docs/stories/importing-and-exporting#exporting-stories) and place the exported filed in the same directory as your `main.tf` file. After that, you can define your story as a Terraform Resource using the following syntax:

```terraform
# provider.tf
provider "tines" {}

# main.tf
resource "tines_story" "dev_story_name" {
  data            = file("${path.module}/story-example.json")
  tenant_url      = "https://dev-tenant.tines.com"
  tines_api_token = var.dev_tines_api_token
  team_id         = var.team_id   # optional
  folder_id       = var.folder_id # optional
}

# variable.tf
variable "dev_tines_api_token" {
  type = string
}

variable "team_id" {
  type = number
}

variable "folder_id" {
  type = number
}
```

And that's all. You don't need to import the state into Terraform either. Running a `terraform apply` will automatically perform an upsert and set the state in Terraform accordingly.