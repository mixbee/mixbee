#!/usr/bin/env bash


# Exit on first error, print all commands.
set -ev

# Shut down the Docker containers that might be currently running.
docker-compose -f docker-compose.yml stop

rm -rf ./node1/Chain ./node1/Log ./node1/peers.recent
rm -rf ./node2/Chain ./node2/Log ./node2/peers.recent
rm -rf ./node3/Chain ./node3/Log ./node3/peers.recent
rm -rf ./node4/Chain ./node4/Log ./node4/peers.recent
rm -rf ./node5/Chain ./node5/Log ./node5/peers.recent
rm -rf ./node6/Chain ./node6/Log ./node6/peers.recent
rm -rf ./node7/Chain ./node7/Log ./node7/peers.recent
rm -rf ./node8/Chain ./node8/Log ./node8/peers.recent


docker rm -f  $(docker ps -a -q)


