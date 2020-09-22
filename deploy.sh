#!/bin/bash
make

VERSION=$(cat version)
chmod +x terraform-provider-metabase_${VERSION}
mv terraform-provider-metabase_${VERSION} ~/.terraform.d/plugins/perxtech.com/tf/metabase/${VERSION}/darwin_amd64/terraform-provider-metabase
rm -rf .terraform/plugins
