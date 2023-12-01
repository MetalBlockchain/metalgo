#!/usr/bin/env bash

set -euo pipefail

###
# cleanup removes the docker instance and the network
echo "Cleaning up..."
docker rm $(sudo docker stop $(sudo docker ps -a -q --filter ancestor=metalblockchain/metalgo:latest --format="{{.ID}}"))  #if the filter returns nothing the command fails, so ignore errors
docker network rm controlled-net 
rm /opt/mainnet-db-daily* 2>/dev/null 
rm -rf /var/lib/metalgo 2>/dev/null 
echo "Done cleaning up"
