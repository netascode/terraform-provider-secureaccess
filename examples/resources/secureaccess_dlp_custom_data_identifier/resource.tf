resource "secureaccess_dlp_custom_data_identifier" "example" {
  name             = "my_dlp_custom_data_identifier"
  description      = "This is my DLP data classification."
  unique_threshold = true
  threshold_ranges = [
    {
      minimum  = 1
      maximum  = 10
      severity = "low"
    }
  ]
  terms     = ["hello_term"]
  patterns  = ["pattern"]
  proximity = ["hello_proximity"]
}
