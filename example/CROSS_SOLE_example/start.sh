#!/usr/bin/env bash


set -ev

### account password is test

docker-compose -f docker-compose.yml down

echo "start main chain(netId=4) node1.example.com......."
docker-compose -f docker-compose.yml up -d node1.example.com

sleep 10

echo "start sub chain(netId=5) node2.example.com......."
docker-compose -f docker-compose.yml up -d node2.example.com

sleep 10
echo "query main chain addr(ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo) balance info when before cross chain"
docker exec node1.example.com mixbee asset balance ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo

echo "query main chain addr(Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR) balance info when before cross chain"
docker exec node1.example.com mixbee asset balance Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR

echo "query sub chain addr(ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo) balance info when before cross chain"
docker exec node2.example.com mixbee asset balance ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo

echo "query sub chain addr(Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR) balance info when before cross chain"
docker exec node2.example.com mixbee asset balance Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR

sleep 20
### 跨链交易发送
echo "main chain addr(Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR) cross chain transfer"
docker exec node1.example.com mixbee asset ctransfer --from ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo --to Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR --aamount 1000 --bamount 10000 --achainid 4 --bchainid 5 --nonce 0 --cpbk 02e750714ff8dfa402a9b7cce6db90de8a72b7e1efe0836e0f65f92b97924e9c7e --password test

echo "sub chain addr(Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR) cross chain transfer"
docker exec node2.example.com mixbee asset ctransfer --from Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR --to ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo --aamount 10000 --bamount 1000 --achainid 5 --bchainid 4 --nonce 0 --cpbk 02e750714ff8dfa402a9b7cce6db90de8a72b7e1efe0836e0f65f92b97924e9c7e --password test

sleep 30
echo "query main chain addr(ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo) balance info when after cross chain"
docker exec node1.example.com mixbee asset balance ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo

echo "query main chain addr(Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR) balance info when after cross chain"
docker exec node1.example.com mixbee asset balance Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR

echo "query sub chain addr(ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo) balance info when after cross chain"
docker exec node2.example.com mixbee asset balance ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo

echo "query sub chain addr(Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR) balance info when after cross chain"
docker exec node2.example.com mixbee asset balance Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR





