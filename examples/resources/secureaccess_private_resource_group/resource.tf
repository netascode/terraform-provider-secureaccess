resource "secureaccess_private_resource_group" "example" {
  name                 = "my private resource group"
  description          = "This is my private resource group"
  private_resource_ids = [123456]
}
