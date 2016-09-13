#!/bin/bash

#################################################
# Get slave ip from docker container for newman
#
# Author: Taka Wang
# Date: 2016/09/13
#################################################

SLAVE=$(docker inspect --format '{{ .NetworkSettings.IPAddress }}' mbweb_slave_1)

echo '{"name": "Docker","values": [{"key": "slave","value": "'"$SLAVE"'","type": "text","enabled": true}]}' > /tmp/postman-env.json