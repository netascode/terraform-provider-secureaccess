resource "secureaccess_zta_profile_private_access_steering_destination" "example" {
  zta_profile_id = "acde070d-8c4c-4f0d-9d8a-162843c10333"
  destination    = "internal.example.com"
  exclusions     = ["app.internal.example.com"]
}
