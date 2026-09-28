resource "secureaccess_network_object_group" "example" {
  name        = "my_network_object"
  description = "This is my network object group."
  network_object_groups = [
    {
      id = 123456
    }
  ]
  network_objects = [
    {
      id = 123456
    }
  ]
  literals = [
    {
      type  = "network"
      value = "192.168.2.0/24"
    }
  ]
}
