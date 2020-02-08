#!/usr/bin/env sh
set -e

terraform init
terraform apply -var ip_address=$(curl ifconfig.co) -var domain=mygaragedoor.space -var a_record_domain=open.mygaragedoor.space -auto-approve