resource "secureaccess_dlp_custom_classification" "example" {
  name                = "my_dlp_custom_classification"
  description         = "My custom classification for sensitive data."
  evaluation_strategy = "any"
  included_data_identifiers = [
    {
      id = "c9023098-7717-40e5-8e94-1f0687f125f6"
    }
  ]
  excluded_data_identifiers = [
    {
      id = "c9023098-7717-40e5-8e94-1f0687f125f6"
    }
  ]
}
