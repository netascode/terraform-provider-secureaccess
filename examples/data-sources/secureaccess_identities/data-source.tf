data "secureaccess_identities" "example" {
  identity_type = "directory_group"
  items = {
    "Finance" = {
    }
  }
}
