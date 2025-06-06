---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "adverity_destination Resource - adverity"
subcategory: ""
description: |-
  Manages a destination.
---

# adverity_destination (Resource)

Manages a destination.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `destination_type_id` (Number) Numeric identifier of the destination type.
- `name` (String) Name of the destination.

### Optional

- `auth_id` (Number) Numeric identifier of the authentication.
- `parameters` (Dynamic) Additional destination parameters.
- `stack_id` (Number) Numeric identifier of the workspace.

### Read-Only

- `id` (Number) Numeric identifier of the destination.
- `last_updated` (String) Timestamp of the last Terraform update of the destination.
