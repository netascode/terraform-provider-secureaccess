resource "secureaccess_network_tunnel_group" "example" {
  name                       = "my network tunnel group"
  device_type                = "ASA"
  region                     = "us-east-1"
  backup_region              = "us-east-2"
  tunnel_id_format_email     = "my_prefix"
  preshared_key              = "t3stingTunn3lNow"
  routing_type               = "static"
  static_network_cidrs       = ["192.168.0.0/24"]
  bgp_override_summarization = true
  ipv6_routing               = true
}
