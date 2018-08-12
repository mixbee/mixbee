#!/usr/bin/env bash

set -ev

echo "query main chain addr(AX5z57vTji4HMCY4N9Ve6Wf16LV2mMoKuT) balance info before cross chain"
docker exec net1node1.example.com mixbee asset balance AX5z57vTji4HMCY4N9Ve6Wf16LV2mMoKuT

echo "query main chain addr(AXwH6juiRvh9G5xVffkA8sczqQFvnySj4f) balance info before cross chain"
docker exec net1node1.example.com mixbee asset balance AXwH6juiRvh9G5xVffkA8sczqQFvnySj4f

echo "query sub chain addr(AX5z57vTji4HMCY4N9Ve6Wf16LV2mMoKuT) balance info  before cross chain"
docker exec net2node2.example.com mixbee asset balance AX5z57vTji4HMCY4N9Ve6Wf16LV2mMoKuT

echo "query sub chain addr(AXwH6juiRvh9G5xVffkA8sczqQFvnySj4f) balance info  before cross chain"
docker exec net2node2.example.com mixbee asset balance AXwH6juiRvh9G5xVffkA8sczqQFvnySj4f


### 跨链交易发送
echo "main chain addr(AX5z57vTji4HMCY4N9Ve6Wf16LV2mMoKuT) cross chain transfer"
docker exec net1node3.example.com mixbee asset ctransfer --from AX5z57vTji4HMCY4N9Ve6Wf16LV2mMoKuT --to AXwH6juiRvh9G5xVffkA8sczqQFvnySj4f --aamount 100 --bamount 1000 --achainid 1 --bchainid 2 --nonce 1 --cpbk 02af0eafe6835fe2bbdfe0c68363b0426db71e0712a22142458c8d82d54c3f7e59 --password test

echo "sub chain addr(Ac6Ac6w3XdhrZWSMvzFyf6WSAtUbBHiFmR) cross chain transfer"
docker exec net2node2.example.com mixbee asset ctransfer --from AXwH6juiRvh9G5xVffkA8sczqQFvnySj4f --to AX5z57vTji4HMCY4N9Ve6Wf16LV2mMoKuT --aamount 1000 --bamount 100 --achainid 2 --bchainid 1 --nonce 1 --cpbk 022d9d9601a7d37b6bd394dfe645e57a891bafba14f1d35f5f5afcb4a58b3c26e8 --password test

sleep 40
echo "query main chain addr(AX5z57vTji4HMCY4N9Ve6Wf16LV2mMoKuT) balance info  before cross chain"
docker exec net1node1.example.com mixbee asset balance AX5z57vTji4HMCY4N9Ve6Wf16LV2mMoKuT

echo "query main chain addr(AXwH6juiRvh9G5xVffkA8sczqQFvnySj4f) balance info  before cross chain"
docker exec net1node1.example.com mixbee asset balance AXwH6juiRvh9G5xVffkA8sczqQFvnySj4f

echo "query sub chain addr(AX5z57vTji4HMCY4N9Ve6Wf16LV2mMoKuT) balance info  before cross chain"
docker exec net2node2.example.com mixbee asset balance AX5z57vTji4HMCY4N9Ve6Wf16LV2mMoKuT

echo "query sub chain addr(AXwH6juiRvh9G5xVffkA8sczqQFvnySj4f) balance info before cross chain"
docker exec net2node2.example.com mixbee asset balance AXwH6juiRvh9G5xVffkA8sczqQFvnySj4f