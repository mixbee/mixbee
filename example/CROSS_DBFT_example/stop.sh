#!/usr/bin/env bash


# Exit on first error, print all commands.
set -ev

# Shut down the Docker containers that might be currently running.
docker-compose -f docker-compose.yml stop

rm -rf ./net1node1/Chain ./net1node1/Log ./net1node1/peers.recent
rm -rf ./net1node2/Chain ./net1node2/Log ./net1node2/peers.recent
rm -rf ./net1node3/Chain ./net1node3/Log ./net1node3/peers.recent
rm -rf ./net1node4/Chain ./net1node4/Log ./net1node4/peers.recent
rm -rf ./net1node5/Chain ./net1node5/Log ./net1node5/peers.recent

rm -rf ./net2node1/Chain ./net2node1/Log ./net2node1/peers.recent
rm -rf ./net2node2/Chain ./net2node2/Log ./net2node2/peers.recent
rm -rf ./net2node3/Chain ./net2node3/Log ./net2node3/peers.recent
rm -rf ./net2node4/Chain ./net2node4/Log ./net2node4/peers.recent
rm -rf ./net2node5/Chain ./net2node5/Log ./net2node5/peers.recent

docker rm -f  $(docker ps -a -q)


