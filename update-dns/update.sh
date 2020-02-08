#!/usr/bin/env sh
set -e

apk add curl
git clone https://github.com/benjefferies/gogo-garage-opener.git
pushd gogo-garage-opener/update-dns
terraform init
terraform apply -var ip_address=$(curl ifconfig.co) -var domain=mygaragedoor.space -var a_record_domain=open.mygaragedoor.space -auto-approve
popd
rm -r gogo-garage-opener