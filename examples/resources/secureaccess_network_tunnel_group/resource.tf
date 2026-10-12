resource "secureaccess_network_tunnel_group" "example" {
  name                   = "my network tunnel group"
  device_type            = "ASA"
  region                 = "us-east-1"
  tunnel_id_format_email = "my_prefix"
  preshared_key          = "t3stingTunn3lNow"
  routing_type           = "static"
  routing_network_cidrs  = ["192.168.0.0/24"]
}
