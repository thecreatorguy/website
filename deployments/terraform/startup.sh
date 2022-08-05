#!/bin/bash

su -

hostnamectl set-hostname itstimjohnson.com

sleep 300 # wait for the EIP to be assigned to this instance so the certbot can work
certbot certonly --standalone --non-interactive --agree-tos -m tim@itstimjohnson.com -d itstimjohnson.com
(crontab -l 2>/dev/null; echo "0 0 1 * * /home/ubuntu/renewcert.sh") | crontab -

su - ubuntu
cd

git clone git@github.com:thecreatorguy/website.git
git checkout add-terraform
cd website
sudo make runprod
sudo make initdb