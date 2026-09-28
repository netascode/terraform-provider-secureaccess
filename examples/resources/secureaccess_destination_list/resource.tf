resource "secureaccess_destination_list" "example" {
  name = "my_destination_list"
  destinations = [
    {
      destination = "example.com"
      comment     = "This is a comment for the destination."
      type        = "domain"
    }
  ]
}
