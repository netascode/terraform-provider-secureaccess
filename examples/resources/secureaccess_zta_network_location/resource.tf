resource "secureaccess_zta_network_location" "example" {
  name = "my_network_location"
  egress_criteria = [
    {
      source_type    = "manual"
      manual_address = "203.0.113.10"
    }
  ]
  egress_exclusion_criteria = [
    {
      source_type    = "manual"
      manual_address = "203.0.113.200"
    }
  ]
  private_subnet_criteria = [
    {
      source_type    = "manual"
      manual_address = "10.10.0.0/16"
    }
  ]
}
