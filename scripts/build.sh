#!/bin/bash

VERSION=$(cat version)
echo "building terraform-provider-metabase_${VERSION}"
go build -o terraform-provider-metabase_${VERSION}