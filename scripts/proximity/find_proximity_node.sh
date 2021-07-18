#!/bin/bash
OFFSET_PORT=21000

NODE_ID=$1
BASE_NODE_SWARM_ADDRESS=$2

#param $1 node id
function load_bee_node_port() {
  OFFSET_NODE=$(expr 3 \* ${NODE_ID})
  API_PORT=$(expr ${OFFSET_PORT} + ${OFFSET_NODE})
  DEBUG_PORT=$(expr ${API_PORT} + 2)
}

function start_one_node() {
  mkdir -p ./nodes/${NODE_ID}
  cp -f ./base/start-template.sh ./nodes/${NODE_ID}/bee${NODE_ID}-private-network.sh
  cd ./nodes/${NODE_ID}

  sed -i 's/<START_PORT>/'"${API_PORT}"'/g' bee${NODE_ID}-private-network.sh
  pm2 -s start bee${NODE_ID}-private-network.sh
  sleep 2
  cd ../../
}

function stop_one_node() {
  pm2 -s stop bee${NODE_ID}-private-network
  pm2 -s delete bee${NODE_ID}-private-network
}

#param $1 node ethereum address
function transfer_token_to_new_node() {
  echo { "Address": ["${NODE_SWARM_ADDRESS}"] } > address.json
  mv address.json ./../code/earthledger/bzzaar-contracts/
  cd ./../code/earthledger/bzzaar-contracts/
  ./node_modules/.bin/yarn transfer:private
  cd ../../../bee_node/
}

function get_node_ethereum_address() {
  NODE_ETHEREUM_ADDRESS=""
  waiting_and_get_node_address "ethereum"
}

function get_node_swarm_address() {
  NODE_SWARM_ADDRESS=""
  waiting_and_get_node_address "swarm"
}

function and_get_node_address() {
  FLAG=$1

  FIND_ADDRESS="false"
  for ((i = 0; i < 3000; i = i + 1)); do
    #    nodeAddress=$(curl -s -X GET  http://localhost:${DEBUG_PORT}/addresses | awk '{split($0, a, "\""); print a[8]}')
    if [ "${FLAG}" == "ethereum" ]; then
      nodeAddress=$(curl -s -X GET http://localhost:${DEBUG_PORT}/addresses | awk '{s=index($0, "ethereum"); print "", substr($0,s+10,45)}')
      if [ -z "${nodeAddress}" ]; then
        sleep 1
      else
        FIND_ADDRESS="true"
        break
      fi
    elif [ "${FLAG}" == "swarm" ]; then
      #      nodeAddress=$(curl -s -X GET http://localhost:${DEBUG_PORT}/addresses | awk '{s=index($0, "/ip4/"); print "", substr($0,s-80,67)}')
      nodeAddress=$(curl -s -X GET http://localhost:${DEBUG_PORT}/addresses)
      if [[ $nodeAddress == *"\"overlay\":null"* ]]; then
        sleep 2
      else
        nodeAddress=$(echo $nodeAddress | awk '{s=index($0, "overlay"); print "", substr($0,s+9,67)}')
        FIND_ADDRESS="true"
        break
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

  echo "${FLAG} ${nodeAddress}"
  if [ "${FLAG}" == "ethereum" ]; then
    NODE_ETHEREUM_ADDRESS=$nodeAddress
  else
    NODE_SWARM_ADDRESS=$nodeAddress
  fi
}

function calc_node_proximity() {
  cd ./../code/earthledger/bee/
  ./dist/tool/proximity $BASE_NODE_SWARM_ADDRESS $NODE_SWARM_ADDRESS
  cd ../../../bee_node/
}

load_bee_node_port $1
start_one_node
get_node_ethereum_address
transfer_token_to_new_node
get_node_swarm_address
calc_node_proximity

#
#for i in $(seq $START $END); do
#  start_one_node $i
#done
