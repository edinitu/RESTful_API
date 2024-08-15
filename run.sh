#/bin/bash

BALANCER_PORT=$1
NUMBER_OF_INSTANCES=$2

rm server/logs_instance_*
rm balancer/logs_load_balancer.txt

cd server
go build .

pids=()

for ((i = 0 ; i < NUMBER_OF_INSTANCES ; i++)); do
  API_PORT=$(($BALANCER_PORT + $i + 1))
  INSTANCE=$(($i+1))
  echo "Will start an api instance on port ${API_PORT}"
  ./server -port=${API_PORT} -no_of_instances=$NUMBER_OF_INSTANCES &>> logs_instance_${INSTANCE}.txt &
  pids+=($!)
done

cd ..
cd balancer
go build .
echo "Will start load balancer on port ${BALANCER_PORT}"
./balancer -number_of_instances=$NUMBER_OF_INSTANCES &>> logs_load_balancer.txt & pids+=($!)

for pid in "${pids[@]}"; do
    wait $pid
done


