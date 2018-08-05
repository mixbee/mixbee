#!/usr/bin/env bash


# Exit on first error, print all commands.
set -ev

# Shut down the Docker containers that might be currently running.
docker-compose -f docker-compose.yml stop

rm -rf ./net1/Chain ./net1/Log ./net1/peers.recent
rm -rf ./net2/Chain ./net2/Log ./net2/peers.recent

docker rm -f  $(docker ps -a -q)


