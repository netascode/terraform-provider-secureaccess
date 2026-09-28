resource "secureaccess_network" "example" {
  name          = "my_network"
  status        = "OPEN"
  is_dynamic    = true
  prefix_length = 30
  prefix        = "192.168.1.0"
}
