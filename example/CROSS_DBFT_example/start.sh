#!/usr/bin/env bash


set -ev

### account1,account2,account3,account4  password is test


docker-compose -f docker-compose.yml down

docker-compose -f docker-compose.yml up -d net1node1.example.com

sleep 5

docker-compose -f docker-compose.yml up -d net1node2.example.com net1node3.example.com net1node4.example.com net1node5.example.com

docker-compose -f docker-compose.yml up -d net2node1.example.com net2node2.example.com net2node3.example.com net2node4.example.com net2node5.example.com
