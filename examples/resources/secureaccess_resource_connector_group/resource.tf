resource "secureaccess_resource_connector_group" "example" {
  name        = "my resource connector group"
  region      = "us-west-2"
  environment = "container"
}
