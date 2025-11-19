# Search for BigQuery authorization
data "adverity_authorization_type" "bigquery" {
  search_term = "bigquery"
}

output "all" {
  value = data.adverity_authorization_type.bigquery.results
}

output "service_account" {
  value = data.adverity_authorization_type.bigquery.results.google-bigquery-service-account
}
