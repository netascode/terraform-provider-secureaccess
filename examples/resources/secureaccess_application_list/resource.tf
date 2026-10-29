resource "secureaccess_application_list" "example" {
  name                     = "my_application_list"
  application_ids          = [1]
  application_category_ids = [40]
  is_default               = false
}
