#!/bin/bash

output_file="queryOrigin.txt"

for i in {1..100}
do
  start_time=$(gdate +%s%N)

  peer chaincode query -C used-car -n used-car-transfer -c '{"function":"ReadTransaction","Args":["1"]}'

  end_time=$(gdate +%s%N)
  elapsed=$(( (end_time - start_time) / 1000 ))


  echo "$elapsed" >> $output_file
   sleep 4
done
