#!/bin/bash

sudo apt update
sudo apt upgrade -y
sudo apt install -y \
    docker.io \
    docker-compose \
    postgresql-client-common \
    make \
    certbot