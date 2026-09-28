resource "secureaccess_internal_network" "example" {
  name          = "my_internal_network"
  prefix        = "10.1.1.0"
  prefix_length = 24
  site_id       = 1819818
}
