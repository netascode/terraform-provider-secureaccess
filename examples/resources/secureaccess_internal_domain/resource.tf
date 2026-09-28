resource "secureaccess_internal_domain" "example" {
  domain                         = "mydomain.local"
  description                    = "My Internal Domain"
  include_all_mobile_devices     = false
  include_all_virtual_appliances = false
  site_ids                       = [123]
}
