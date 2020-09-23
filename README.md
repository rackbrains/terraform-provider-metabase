# terraform-provider-metabase

Build and deploy locally using

```bash
make
```

Use in terraform template using

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

resource "metabase_card" "nicolas_test" {
  name          = "nicolas test"
  description   = "metabase terraform provider test"
  query         = "select * from foo"
  collection_id = 26
  variables {
    id           = "1"
    name         = "start_date"
    type         = "date"
    display_name = "Start Date"
    required     = true
  }
  variables {
    id           = "2"
    name         = "end_date"
    type         = "date"
    display_name = "End Date"
    required     = true
  }
}
```
