resource "secureaccess_private_resource" "example" {
  name           = "my private resource"
  description    = "This is my private resource"
  certificate_id = 123456
  dns_server_id  = 789012
  endpoint_connection_methods = [
    {
      type                       = "client"
      client_reachable_addresses = ["192.168.1.0/24"]
    }
  ]
  addresses = [
    {
      addresses = ["192.168.1.0/24"]
      protocols_and_ports = [
        {
          protocol = "http/https"
          ports    = "443"
        }
      ]
    }
  ]
}
