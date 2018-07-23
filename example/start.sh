#!/usr/bin/env bash


set -ev

### account1,account2,account3,account4 is test


docker-compose -f docker-compose.yml down


docker-compose -f docker-compose.yml up -d node1.example.com node2.example.com node3.example.com node4.example.com
