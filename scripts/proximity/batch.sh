#!/bin/bash

START=$1
END=$2
BATCH=$3

#  for ((i = 0; i < length; i = i + 4)); do

#for ((i = 6; i < 600; i = i + 100)); do
	#  echo $i $(expr $i + 29)
	#  echo "./get_node_address.sh ${i} $(expr ${i} + 99) swarm > swarm_address_${i}.log &"
	  #  nohup bash -C './get_node_address.sh $i $(expr $i + 29) swarm' > swarm_address_$i.log &
#done


for ((i = 6; i < 600; i = i + 100)); do
    echo swarm_address_${i}.log
	cat swarm_address_${i}.log
done

# pm2 stop $(pm2 status | grep private-network | awk '{print $4}')
# pm2 delete $(pm2 status | grep private-network | awk '{print $4}')
