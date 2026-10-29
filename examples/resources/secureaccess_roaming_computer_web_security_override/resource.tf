resource "secureaccess_roaming_computer_web_security_override" "example" {
  enabled_origin_ids  = [123456789]
  disabled_origin_ids = [987654321]
}
