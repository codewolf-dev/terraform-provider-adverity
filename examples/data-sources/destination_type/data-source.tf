data "adverity_destination_type" "bigquery" {
  search_term = "bigquery"
}

output "all" {
  value = data.adverity_destination_type.bigquery.results
}
