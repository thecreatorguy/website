#!/bin/bash
(
    echo "Renewing ssl certificate `date`"
    cd /home/ubuntu/website
    sudo make stopprod
    certbot certonly --standalone --force-renew -d itstimjohnson.com
    sudo make runprod
    echo ""
) >> ~/cronlog 2>&1