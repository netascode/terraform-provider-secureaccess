resource "secureaccess_do_not_decrypt_list" "example" {
  name = "my_do_not_decrypt_list"
  destinations = [
    {
      destination = "example.com"
      comment     = "This is a comment for the destination."
    }
  ]
  applications                 = [101]
  decrypt_exception_categories = [12]
}
