---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "tines_story Resource - terraform-provider-tines"
subcategory: ""
description: |-
  Manage a Tines Story.
---

# tines_story (Resource)

Manage a Tines Story.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `data` (String) Tines Story export that gets read in from a JSON file
- `tenant_url` (String) Tines tenant URL
- `tines_api_token` (String, Sensitive) API token for Tines Tenant

### Optional

- `folder_id` (Number) Tines folder ID.
- `team_id` (Number) Tines team ID.

### Read-Only

- `id` (Number) Tines Story identifier.
