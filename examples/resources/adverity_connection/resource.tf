# Deprecated - use adverity_authorization instead
resource "adverity_connection" "sprinklr" {
  name     = "sprinklr"
  stack_id = 1

  connection_type_id = 187 # Sprinklr

  parameters = {
    domain = "your-spinklr-space"
  }
}

# Deprecated - use adverity_authorization instead
resource "adverity_connection" "bigquery" {
  name     = "bigquery"
  stack_id = 1

  connection_type_id = 284 # BigQuery

  parameters = {
    base64_encoded_credentials = filebase64("path/to/credentials.json")
  }
}
