#!/usr/bin/env sh
set -e

curl -L -o /tmp/terraform.zip https://releases.hashicorp.com/terraform/0.12.24/terraform_0.12.24_linux_arm.zip
unzip -o /tmp/terraform.zip -d /tmp
chmod +x /tmp/terraform
/tmp/terraform init
/tmp/terraform apply -var ip_address=$(curl ifconfig.co) -var domain=mygaragedoor.space -var a_record_domain=open.mygaragedoor.space -auto-approve