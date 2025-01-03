#!/bin/bash

kafka-topics --bootstrap-server kafka-0:9092 --list

kafka-topics --create --if-not-exists --topic mp-messages-stream --bootstrap-server kafka-0:9092 --partitions 3 --replication-factor 2
kafka-topics --create --if-not-exists --topic mp-blocked-users-stream --bootstrap-server kafka-0:9092 --partitions 3 --replication-factor 2
kafka-topics --create --if-not-exists --topic mp-filtered-messages-stream --bootstrap-server kafka-0:9092 --partitions 3 --replication-factor 2

kafka-topics --create --if-not-exists --topic block-users-consumer-table --bootstrap-server kafka-0:9092 --partitions 3 --replication-factor 2 --config cleanup.policy=compact