resource "secureaccess_zta_profile" "example" {
  name                                    = "my ztna profile"
  priority                                = 1
  operating_system_macos                  = true
  operating_system_windows                = true
  operating_system_linux                  = false
  operating_system_ios                    = true
  operating_system_android                = true
  operating_system_android_knox           = false
  all_users                               = true
  all_groups                              = true
  network_location_ids                    = ["8f0831d9-02e2-4c2d-9f71-47874e9465f7"]
  private_access_use_trusted_networks     = true
  private_access_end_user_pause           = true
  private_access_end_user_pause_duration  = 30
  internet_access_steering_mode           = 1
  internet_access_use_trusted_networks    = true
  internet_access_end_user_pause          = true
  internet_access_end_user_pause_duration = 60
}
