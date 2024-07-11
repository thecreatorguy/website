#!/bin/bash
hostnamectl set-hostname itstimjohnson.com

sleep 5 # wait for the EIP to be assigned to this instance so the certbot can work
certbot certonly --standalone --non-interactive --agree-tos -m tim@itstimjohnson.com -d itstimjohnson.com
(crontab -l 2>/dev/null; echo "0 0 1 * * /home/ubuntu/renewcert.sh") | crontab -

sudo -i -u ubuntu ./start_website.sh