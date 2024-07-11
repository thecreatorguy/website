#!/bin/bash
cd

git clone https://github.com/thecreatorguy/website.git
cd website
aws secretsmanager get-secret-value --secret-id website/env --output=json --no-cli-pager | jq -r .SecretString > ./build/package/.env
sudo make initmodules
sudo make runprod
sudo make initdb