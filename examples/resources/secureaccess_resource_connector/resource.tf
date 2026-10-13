resource "secureaccess_resource_connector" "example" {
  instance_id       = "i-0123456789abdef1"
  enabled           = true
  confirmed         = true
  delete_on_destroy = false
}
