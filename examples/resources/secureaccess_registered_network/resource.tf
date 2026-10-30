resource "secureaccess_registered_network" "example" {
  name               = "my_network"
  status             = "OPEN"
  is_dynamic         = true
  ipv4_prefix        = "198.51.100.0"
  ipv4_prefix_length = 32
}
