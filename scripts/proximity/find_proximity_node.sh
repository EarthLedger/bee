#!/bin/bash
OFFSET_PORT=21000

START=$1
END=$2
BASE_NODE_SWARM_ADDRESS=$3
TARGET_PROXIMITY=$4

#param $1 node id
function load_bee_node_port() {
  OFFSET_NODE=$(expr 3 \* ${NODE_ID})
  API_PORT=$(expr ${OFFSET_PORT} + ${OFFSET_NODE})
  DEBUG_PORT=$(expr ${API_PORT} + 2)
}

function start_one_node() {
  echo "... start node ..."
  mkdir -p ./nodes/${NODE_ID}
  cp -f ./base/start-template.sh ./nodes/${NODE_ID}/bee${NODE_ID}-private-network.sh
  cd ./nodes/${NODE_ID}

  sed -i 's/<START_PORT>/'"${API_PORT}"'/g' bee${NODE_ID}-private-network.sh
  pm2 -s start bee${NODE_ID}-private-network.sh
  sleep 2
  cd ../../
}

function stop_one_node() {
  echo "... stop node ..."

  pm2 -s stop bee${NODE_ID}-private-network
  pm2 -s delete bee${NODE_ID}-private-network
}

#param $1 node ethereum address
function transfer_token_to_new_node() {
  echo "... trasfer token  ..."
  echo --${NODE_ETHEREUM_ADDRESS}--
  cp address-template.json address.json
  sed -i 's/<NODE_ADDRESS>/'"${NODE_ETHEREUM_ADDRESS}"'/g' address.json
  mv address.json ./../code/earthledger/bzzaar-contracts/
  cd ./../code/earthledger/bzzaar-contracts/
  ./node_modules/.bin/yarn transfer:private
  cd ../../../bee_node/
}

function get_node_ethereum_address() {
  echo "... get node ethereum address  ..."

  NODE_ETHEREUM_ADDRESS=""
  get_node_address "ethereum"
}

function get_node_swarm_address() {
  echo "... get node swarm address  ..."

  NODE_SWARM_ADDRESS=""
  get_node_address "swarm"
}

function get_node_address() {
  FLAG=$1

  FIND_ADDRESS="false"
  for ((i = 0; i < 3000; i = i + 1)); do
    #    nodeAddress=$(curl -s -X GET  http://localhost:${DEBUG_PORT}/addresses | awk '{split($0, a, "\""); print a[8]}')
    if [ "${FLAG}" == "ethereum" ]; then
      nodeAddress=$(curl -s -X GET http://localhost:${DEBUG_PORT}/addresses | awk '{s=index($0, "ethereum"); print  substr($0,s+11,42)}')
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
        nodeAddress=$(echo $nodeAddress | awk '{s=index($0, "overlay"); print  substr($0,s+10,64)}')
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

  #  echo "${FLAG} ${nodeAddress}"
  if [ "${FLAG}" == "ethereum" ]; then
    NODE_ETHEREUM_ADDRESS=$nodeAddress
  else
    NODE_SWARM_ADDRESS=$nodeAddress
  fi
}

function calc_node_proximity() {
  echo "... calc proximity  ..."
  echo --$BASE_NODE_SWARM_ADDRESS--:--$NODE_SWARM_ADDRESS--
  cd ./../code/earthledger/bee/
  ./dist/tool/proximity $BASE_NODE_SWARM_ADDRESS $NODE_SWARM_ADDRESS >proximity.log
  node_proximity_str=$(cat proximity.log | awk '{s=index($0, "proximity"); print substr($0,s+11,2)}')
  if [ -z "${node_proximity_str}" ]; then
    echo "do proximity calc error"
    exit 255
  else
    node_proximity=$(expr $node_proximity_str)
    if [ $node_proximity -ge $TARGET_PROXIMITY ]; then
      FIND_TARGET="true"
    fi
  fi

  cd ../../../bee_node/
}

function proximity_one_node() {
  load_bee_node_port
  start_one_node
  get_node_ethereum_address
  transfer_token_to_new_node
  get_node_swarm_address
  calc_node_proximity
  stop_one_node
}

function calc_run_time() {
  END_TIME=$(date +%s)
  DURATION=$(( $END_TIME - $START_TIME ))
  DURATION_M=$(expr $DURATION / 60)
  echo "total time: $DURATION second / $DURATION_M minute   "
}

START_TIME=$(date +%s)

for i in $(seq $START $END); do
  NODE_ID=$1
  proximity_one_node
  if [ "${FIND_TARGET}" == "true" ]; then
    count=$(expr ${END} - ${START} + 1)
    echo "success, target proximity: ${TARGET_PROXIMITY}, node count:  ${count}"
    calc_run_time
    exit 1
  fi
done

count=$(expr ${END} - ${START} + 1)
echo "failed, target proximity: ${TARGET_PROXIMITY} node count:  ${count}"
calc_run_time
