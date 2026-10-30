resource "secureaccess_access_rule" "example" {
  name                                 = "my_internal_network"
  description                          = "Allows access to internal network"
  action                               = "allow"
  priority                             = 1
  enabled                              = true
  source_identity_ids                  = [683051912]
  source_identity_type_ids             = [7]
  source_network_object_ids            = [501078608]
  source_private_literals              = ["192.168.0.1"]
  source_operator                      = "OR"
  destination_network_object_ids       = [501078608]
  destination_network_object_group_ids = [500018294]
  destination_service_object_ids       = [501078537]
  destination_service_object_group_ids = [500018292]
  destination_operator                 = "AND"
  destination_operator_network_service = "AND"
  log_level                            = "LOG_ALL"
  traffic_type                         = "PRIVATE_NETWORK"
  public_security_profile_id           = 15256616
  private_security_profile_id          = 15256616
  tenant_control_profile_id            = 1281818
  ips_profile_id                       = 524837
}
