resource "secureaccess_zta_profile_internet_access_exceptions" "example" {
  zta_profile_id = "acde070d-8c4c-4f0d-9d8a-162843c10333"
  exceptions = [
    {
      destination = "*.example.com"
      description = "My exception"
    }
  ]
}
