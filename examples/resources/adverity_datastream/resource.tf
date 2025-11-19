resource "adverity_authorization" "sprinklr" {
  name     = "sprinklr"
  stack_id = 1

  connection_type_id = 187 # Sprinklr

  parameters = {
    domain = "your-spinklr-space"
  }
}

resource "adverity_datastream" "datastream" {
  name     = "sprinklr"
  auth_id  = adverity_authorization.sprinklr.id
  stack_id = 1

  datastream_type_id = 576 # Sprinklr

  datatype = "Live" # Either "Live" to enable scheduling or "Staging" to disable scheduling

  enabled = false # Enable data transfers to destination

  parameters = {
    widget_query = jsondecode(file("path/to/file.json")) # Pass parameters as Terraform types instead of string
  }

  schedule {
    cron_preset       = "CRON_EVERY_DAY"
    cron_start_of_day = "03:33:33"
    time_range_preset = 0 # Custom
    fixed_start       = "2025-01-01"
  }
}
