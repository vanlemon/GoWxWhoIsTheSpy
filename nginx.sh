#! /bin/bash

cd $(dirname $0)

sudo mv ./nginx.conf /etc/nginx/

sudo nginx -t
sudo nginx -s reload