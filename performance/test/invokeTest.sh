#!/bin/bash

output_file="invoke.txt"

for i in {1..100}
do
  start_time=$(gdate +%s%N)

  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls true --cafile $ORDERER_CA -C used-car -n transaction --peerAddresses localhost:7051 --tlsRootCertFiles ${CORE_PEER_TLS_ROOTCERT_FILE} --peerAddresses localhost:9051 --tlsRootCertFiles /Users/jeho/go/src/github.com/hyperledger/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"SellVehicle", "Args" : ["{\"name\":\"John\",\"residentRegistrationNumber\":\"880808-1******\",\"phoneNumber\":\"123-4567-8910\",\"address\":\"Seoul\"}", "{\"transactionState\":\"\",\"vehicleRegistrationNumber\":\"123A4567\",\"newVehicleRegistrationNumber\":\"\",\"vehicleModelName\":\"Tesla Model S\",\"vehicleIdentificationNumber\":\"5YJ3E1EA1JF00001\",\"transactionDate\":\"\",\"transactionAmount\":\"30000\",\"balancePaymentDate\":\"\",\"vehicleDeliveryDate\":\"\",\"vehicleDeliveryAddress\":\"\",\"mileage\":\"10000\"}"]}'

  end_time=$(gdate +%s%N)

  echo -e "$start_time $end_time" >> $output_file
  sleep 4
done
