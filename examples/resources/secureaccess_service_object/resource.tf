resource "secureaccess_service_object" "example" {
  name        = "my_service_object"
  description = "This is my service object."
  protocol    = "tcp"
  ports       = ["8080"]
}
