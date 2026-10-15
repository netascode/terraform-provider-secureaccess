resource "secureaccess_service_object_group" "example" {
  name        = "my_service_object_group"
  description = "This is my service object group."
  service_object_groups = [
    {
      id = 123456
    }
  ]
  service_objects = [
    {
      id = 123456
    }
  ]
  literals = [
    {
      protocol = "tcp"
      ports    = ["8443"]
    }
  ]
}
