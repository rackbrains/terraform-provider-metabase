terraform {
  required_providers {
    metabase = {
      versions = ["0.1.0"]
      source   = "perxtech.com/tf/metabase"
    }
  }
}

provider "metabase" {
  host     = ""
  username = "..."
  password = "..."
}

# resource "metabase_card" "nicolas_test" {
#   name = "nicolas_test"
# }

# resource "metabase_card" "yohan" {
#   name = "Yohan"
# }