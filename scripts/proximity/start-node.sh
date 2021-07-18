#!/bin/bash
OFFSET_PORT=21000
NODE_ID=$1

mkdir -p ./nodes/${NODE_ID}

cp -f ./base/start-template.sh ./nodes/${NODE_ID}/bee${NODE_ID}-private-network.sh
cd ./nodes/${NODE_ID}

OFFSET_NODE=`expr 3 \* ${NODE_ID}`
API_PORT=`expr ${OFFSET_PORT} + ${OFFSET_NODE}`
sed -i 's/<START_PORT>/'"${API_PORT}"'/g' bee${NODE_ID}-private-network.sh

pm2 start bee${NODE_ID}-private-network.sh
cd ..

