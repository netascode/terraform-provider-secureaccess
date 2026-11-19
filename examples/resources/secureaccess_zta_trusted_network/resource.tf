resource "secureaccess_zta_trusted_network" "example" {
  name    = "my zta trusted network"
  inspect = 1
  dns_servers = [
    {
      ip = "1.1.1.1"
    }
  ]
  dns_domains = [
    {
      name = "*.example.com"
    }
  ]
  trusted_servers = [
    {
      url              = "https://trusted.example.com:8443"
      certificate_hash = "x2v1b7d9p0e3f6g8a4w5z2q1k9m7r5n8l6j0t3u2s9y8v5c1b6h4d7w2x9g3p5a1"
    }
  ]
}
