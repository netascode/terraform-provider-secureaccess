resource "secureaccess_zta_profile_internet_access_destination_list_mappings" "example" {
  zta_profile_id = "acde070d-8c4c-4f0d-9d8a-162843c10333"
  destination_lists = [
    {
      id         = "54321"
      exclusions = ["www.example.com"]
    }
  ]
}
