resource "secureaccess_dlp_ai_guardrails_rule" "example" {
  name                                   = "my_dlp_ai_guardrails_rule"
  description                            = "My AI Guardrails rule"
  severity                               = "CRITICAL"
  enabled                                = true
  data_classification_scannable_contexts = ["CONTENT"]
  data_classifications = [
    {
      id   = "d309239a-b6fd-11ef-a825-0242ac120002"
      name = "Privacy Guardrail"
    }
  ]
  file_labels                  = ["Confidential"]
  file_size_configuration_type = "RANGE"
  file_size_from               = 1024
  file_size_from_unit          = "KB"
  file_size_to                 = 10240
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
  action                         = "BLOCK"
  override_block                 = false
  notification_popup             = true
  email_notify_actor             = true
  email_notify_actors_manager    = true
  email_notify_custom_recipients = ["user1@example.com"]
}
