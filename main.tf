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
