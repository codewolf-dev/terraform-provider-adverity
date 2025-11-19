resource "adverity_connection" "bigquery" {
  name     = "bigquery"
  stack_id = 1

  connection_type_id = 284 # BigQuery

  parameters = {
    base64_encoded_credentials = filebase64("path/to/credentials.json")
  }
}

resource "adverity_destination" "bigquery" {
  name     = "bigquery"
  auth_id  = adverity_connection.bigquery.id
  stack_id = 1

  destination_type_id = 253 # BigQuery

  parameters = {
    schema_mapping     = true
    project            = "example-project"
    dataset            = "example-dataset"
    headers_formatting = 3 # replace spaces by underscores and convert letters to lowercase
  }
}

