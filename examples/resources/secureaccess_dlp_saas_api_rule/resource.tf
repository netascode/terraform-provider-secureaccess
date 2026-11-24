resource "secureaccess_dlp_saas_api_rule" "example" {
  name                                   = "my_saas_api_rule"
  description                            = "My SaaS API rule"
  severity                               = "CRITICAL"
  enabled                                = true
  data_classification_scannable_contexts = ["CONTENT"]
  data_classifications = [
    {
      id   = "39b3b945-2e21-4831-bbbc-86fc5200ddf7"
      name = "Built-in PCI Classification"
    }
  ]
  platform_name                = "spark"
  platform_id                  = 12345
  platform_application_labels  = ["spark"]
  mip_and_titus_labels         = ["Confidential"]
  file_size_configuration_type = "RANGE"
  file_size_from               = 1024
  file_size_from_unit          = "KB"
  file_size_to                 = 51200
  file_size_to_unit            = "KB"
  file_types_included = [
    {
      id   = 250
      name = "cdr"
    }
  ]
  file_type_groups_included = [
    {
      id   = 6
      name = "Documents"
    }
  ]
  file_types_excluded = [
    {
      id   = 156
      name = "dbf"
    }
  ]
  file_type_groups_excluded = [
    {
      id   = 2
      name = "Compressed files"
    }
  ]
  unclassified_files_scan = false
  identities_included = [
    {
      id      = 1442790899
      type_id = 7
    }
  ]
  identities_excluded = [
    {
      type_id = 40
    }
  ]
  user_domain_destinations   = [""]
  file_owner_email_addresses = [""]
  exposures_included = [
    {
      type    = ""
      users   = [""]
      domains = [""]
    }
  ]
  action                         = "MONITOR"
  email_notify_owner             = true
  email_notify_actor             = true
  email_notify_actors_manager    = true
  email_notify_custom_recipients = ["user1@example.com"]
}
