#!/usr/bin/env bash


set -ev

### account1,account2,account3,account4  password is test


docker-compose -f docker-compose.yml down


docker-compose -f docker-compose.yml up -d node1.example.com node2.example.com node3.example.com node4.example.com node5.example.com node6.example.com node7.example.com node8.example.com
