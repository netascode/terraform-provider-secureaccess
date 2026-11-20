resource "secureaccess_zta_profile_group_mappings" "example" {
  zta_profile_id = "acde070d-8c4c-4f0d-9d8a-162843c10333"
  groups = [
    {
      id = "12345678"
    }
  ]
}
