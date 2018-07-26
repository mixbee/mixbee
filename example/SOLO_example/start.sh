#!/usr/bin/env bash


set -ev

### account1,account2,account3,account4  password is test


docker-compose -f docker-compose.yml down


docker-compose -f docker-compose.yml up -d node.example.com