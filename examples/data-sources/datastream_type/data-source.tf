data "adverity_datastream_type" "sprinklr" {
  search_term = "sprinklr"
}

output "all" {
  value = data.adverity_datastream_type.sprinklr.results
}
