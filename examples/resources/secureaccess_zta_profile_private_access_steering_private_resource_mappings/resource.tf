resource "secureaccess_zta_profile_private_access_steering_private_resource_mappings" "example" {
  zta_profile_id = "acde070d-8c4c-4f0d-9d8a-162843c10333"
  private_resources = [
    {
      id = "123456"
    }
  ]
}
