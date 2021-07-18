#!/bin/bash
OFFSET_PORT=21000
START=$1
END=$2
FLAG=$3

function start_one_node() {
  NODE_ID=$1

  mkdir -p ./nodes/${NODE_ID}

  cp -f ./base/start-template.sh ./nodes/${NODE_ID}/bee${NODE_ID}-private-network.sh
  cd ./nodes/${NODE_ID}

  OFFSET_NODE=`expr 3 \* ${NODE_ID}`
  API_PORT=`expr ${OFFSET_PORT} + ${OFFSET_NODE}`
  sed -i 's/<START_PORT>/'"${API_PORT}"'/g' bee${NODE_ID}-private-network.sh

  pm2 -s start bee${NODE_ID}-private-network.sh
  sleep 2

  DEBUG_PORT=`expr ${API_PORT} + 2`
  FIND_ADDRESS="false"
  for ((i = 0; i < 3000; i = i + 1)); do
    #    nodeAddress=$(curl -s -X GET  http://localhost:${DEBUG_PORT}/addresses | awk '{split($0, a, "\""); print a[8]}')
    if [ "${FLAG}" == "ethereum" ]; then
      nodeAddress=$(curl -s -X GET  http://localhost:${DEBUG_PORT}/addresses | awk '{s=index($0, "ethereum"); print "", substr($0,s+10,45)}')
      if [ -z "${nodeAddress}" ]; then
        sleep 1
      else
        FIND_ADDRESS="true"
        break
      fi
    elif [ "${FLAG}" == "swarm" ]; then
#      nodeAddress=$(curl -s -X GET http://localhost:${DEBUG_PORT}/addresses | awk '{s=index($0, "/ip4/"); print "", substr($0,s-80,67)}')
       nodeAddress=$(curl -s -X GET http://localhost:${DEBUG_PORT}/addresses)
       if [[ $nodeAddress == *"/ip4/"* ]]; then
         nodeAddress=$( echo $nodeAddress | awk '{s=index($0, "/ip4/"); print "", substr($0,s-80,67)}')
         FIND_ADDRESS="true"
         break
       else
         sleep 2
       fi
    else
      echo "wrong flag"
      exit 255
    fi
  done

  if [ ${FIND_ADDRESS} == "false" ]; then
    echo "cann't find address from ${NODE_ID}"
    exit 255
  fi

  echo $nodeAddress
  pm2 -s stop bee${NODE_ID}-private-network
  pm2 -s delete bee${NODE_ID}-private-network

  cd ../../
}

for i in $(seq $START $END); do
  start_one_node $i
done

