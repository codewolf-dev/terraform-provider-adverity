resource "adverity_authorization" "sprinklr" {
  name     = "sprinklr"
  stack_id = 1

  connection_type_id = 187 # Sprinklr

  parameters = {
    domain = "your-spinklr-space"
  }
}

resource "adverity_authorization" "bigquery_authorization" {
  name     = "bigquery-connection"
  stack_id = 1

  authorization_type_id = 284 # BigQuery

  parameters = {
    base64_encoded_credentials = filebase64("path/to/credentials.json")
  }
}
