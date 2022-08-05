#!/bin/bash

echo "Renewing ssl certificate `date`" >> ~/cronlog
certbot certonly --standalone --force-renew -d itstimjohnson.com >> ~/cronlog
echo "" >> ~/cronlog