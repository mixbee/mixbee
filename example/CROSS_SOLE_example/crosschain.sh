#!/usr/bin/env bash


set -ev

echo "query main chain addr(ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo) balance info when before cross chain"
docker exec node1.example.com mixbee asset balance ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo

echo "query main chain addr(Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR) balance info when before cross chain"
docker exec node1.example.com mixbee asset balance Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR

echo "query sub chain addr(ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo) balance info when before cross chain"
docker exec node2.example.com mixbee asset balance ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo

echo "query sub chain addr(Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR) balance info when before cross chain"
docker exec node2.example.com mixbee asset balance Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR

### 跨链交易发送
echo "main chain addr(Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR) cross chain transfer"
docker exec node1.example.com mixbee asset ctransfer --from ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo --to Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR --aamount 1000 --bamount 10000 --achainid 4 --bchainid 5 --nonce 1 --cpbk 02e750714ff8dfa402a9b7cce6db90de8a72b7e1efe0836e0f65f92b97924e9c7e --password test

echo "sub chain addr(Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR) cross chain transfer"
docker exec node2.example.com mixbee asset ctransfer --from Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR --to ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo --aamount 10000 --bamount 1000 --achainid 5 --bchainid 4 --nonce 1 --cpbk 02e750714ff8dfa402a9b7cce6db90de8a72b7e1efe0836e0f65f92b97924e9c7e --password test

sleep 50
echo "query main chain addr(ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo) balance info when after cross chain"
docker exec node1.example.com mixbee asset balance ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo

echo "query main chain addr(Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR) balance info when after cross chain"
docker exec node1.example.com mixbee asset balance Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR

echo "query sub chain addr(ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo) balance info when after cross chain"
docker exec node2.example.com mixbee asset balance ANZ7Yh26agCtQjQYAZHZqpXL3jQVPCigxo

echo "query sub chain addr(Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR) balance info when after cross chain"
docker exec node2.example.com mixbee asset balance Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR





