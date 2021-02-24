terraform {
  required_providers {
    metabase = {
      version = "0.3.0"
      source  = "perxtech.com/tf/metabase"
    }
  }
}

provider "metabase" {
  host     = "http://localhost:3000"
  username = "thomasnyambati@gmail.com"
  password = "DUH_.@E3ChubN4$"
}

resource "metabase_card" "test" {
  name             = "Terraform test"
  description      = "metabase terraform provider test"
  query            = "select * from jo"
  collection_id    = 26
  enable_embedding = true
  connection_id    = 15
  variables {
    id              = "1"
    name            = "start_date"
    type            = "date"
    display_name    = "Start Date"
    required        = true
    embedding_param = "enabled"

  }
  variables {
    id              = "2"
    name            = "end_date"
    type            = "date"
    display_name    = "End Date"
    required        = true
    embedding_param = "locked"
  }
}


resource "metabase_database" "name" {
  name     = "sentry"
  engine   = "mysql"
  user     = "value"
  host     = "value"
  db       = ""
  port     = "value"
  password = "value"
}
