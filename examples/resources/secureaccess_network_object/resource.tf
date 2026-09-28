resource "secureaccess_network_object" "example" {
  name        = "my_network_object"
  description = "This is my network object."
  type        = "network"
  value       = "192.168.1.0/24"
}
