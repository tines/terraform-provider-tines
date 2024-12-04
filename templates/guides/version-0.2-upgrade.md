---
layout: ""
page_title: "Upgrading to version 0.2.x (from 0.1.x)"
description: Terraform Tines Provider Version 0.2 Upgrade Guide
---

# Terraform Tines Provider Version 0.2 Upgrade Guide

Version 0.2 has made some architectural and configurability changes under the hood in preparation for some significant new functionality, which will be coming in a future release. While all existing resources are compatible with v0.2, some attributes have been deprecated and some new configuration values are required.

## Provider Version Configuration
If you are not ready to make a move to version 0.2 of the Tines provider, you may keep the 0.1.x branch active for
your Terraform project by specifying:

```terraform
provider "tines" {
  version = "~> 0.1"
  # ... any other configuration
}
```

We highly recommend that you review this guide, make necessary changes and move to 0.2.x branch, as further 0.1.x releases are
unlikely to happen.

~> Before attempting to upgrade to version 0.2, you should first upgrade to the
   latest version of 0.1 to ensure any transitional updates are applied to your
   existing configuration.

## Provider Global Configuration Changes
The following changes have been made at the provider level:

- Added a new required configuration value for `tenant` which can be set either in the provider configuration or as the `TINES_TENANT` environment variable.
- Added a new required configuration value for `api_key` which can be set either in the provider configuration or as the `TINES_API_KEY` environment variable.

## Tines Story Configuration Changes
The following changes have been made to the `tines_story` resource:

- The `tines_api_token` resource attribute has been marked as deprecated and will be removed in a future release. In version 0.2, any value set here will be ignored and overridden by the provider-level `api_key` attribute.
- The `tenant_url` resource attribute has been marked as deprecated and will be removed in a future release. In version 0.2, any value set here will be ignored and overriden by the provider-level `tenant` attibute.