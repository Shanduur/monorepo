#!/bin/env bash

mkfile() { 
    mkdir -p $( dirname "$1") && touch "$1" 
}

mkdir -p /etc/traefik/
mv traefik.yml /etc/traefik/traefik.yml

cd /usr/bin
wget https://github.com/containous/traefik/releases/download/v1.1.2/traefik_linux-amd64 
mv traefik_linux-amd64 traefik 
chmod u+x traefik
cd ~

mkfile /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
mkfile /var/run/secrets/kubernetes.io/serviceaccount/token

mv traefik.service /etc/systemd/system/traefik.service
