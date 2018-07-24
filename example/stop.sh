#!/usr/bin/env bash


# Exit on first error, print all commands.
set -ev

# Shut down the Docker containers that might be currently running.
docker-compose -f docker-compose.yml stop

rm -rf ./config/node1/Chain ./config/node1/Log ./config/node1/peers.recent
rm -rf ./config/node2/Chain ./config/node2/Log ./config/node2/peers.recent
rm -rf ./config/node3/Chain ./config/node3/Log ./config/node3/peers.recent
rm -rf ./config/node4/Chain ./config/node4/Log ./config/node4/peers.recent
rm -rf ./config/node5/Chain ./config/node5/Log ./config/node5/peers.recent

docker rm -f  $(docker ps -a -q)


