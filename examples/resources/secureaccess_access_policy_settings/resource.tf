resource "secureaccess_access_policy_settings" "example" {
  zero_trust_access_authentication         = true
  zero_trust_access_authentication_timeout = 485
  editing_session_duration                 = 55
  editing_session_extension                = 6
  global_ips_profile_id                    = "654321"
  microsoft_365_compatibility              = true
  decryption_global                        = true
  decryption_disable_for_sources           = ["10.10.0.0/24"]
  log_decrypted_internet_destinations      = true
  log_decrypted_private_resources          = true
}
