#!/bin/env bash

firewall-cmd --add-port=3270/tcp
firewall-cmd --runtime-to-permanent
firewall-cmd --reload

yum install -y libnsl unzip
mkdir -p /usr/mvs38 && cd /usr/mvs38

curl -sSLO http://www.hercules-390.org/hercules-3.07-1.x86_64.rpm
curl -sSLO http://wotho.ethz.ch/tk4-/tk4-_v1.00_current.zip
curl -sSLO http://wotho.ethz.ch/tk4-/tk4-source.zip
curl -sSLO http://wotho.ethz.ch/tk4-/tk4-cbt.zip

sudo rpm -i ./hercules-3.07-1.x86_64.rpm

for z in ./*.zip
do
    unzip $z
done
rm *.zip *.rpm

sed -i 's/NUMCPU ${NUMCPU:=1}/NUMCPU ${NUMCPU:=2}/g' ./conf/tk4-.cnf
sed -i 's/MAXCPU ${MAXCPU:=1}/MAXCPU ${MAXCPU:=2}/g' ./conf/tk4-.cnf

mv /tmp/hercules.service /etc/systemd/system/hercules.service

systemctl daemon-reload
systemctl enable hercules.service
systemctl start hercules.service
