# terraform-provider-metabase

## Deploy

Build and deploy locally using

```bash
make
```

## Example usage

```tf
terraform {
  required_providers {
    metabase = {
      versions = ["0.1.0"]
      source   = "perxtech.com/tf/metabase"
    }
  }
}

provider "metabase" {
  host     = "https://metabase.perxtech.io"
  username = "xxxx@perxtech.com"
  password = "yyyy"
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
```

## Import

Import using card id.

```bash
terraform import metabase_card.test 1243
```
