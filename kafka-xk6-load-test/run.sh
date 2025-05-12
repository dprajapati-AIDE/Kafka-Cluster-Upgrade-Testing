#!/bin/bash

SCRIPT_NAME="kafka-load-test.js"
K6_EXECUTABLE="./k6"

export BOOTSTRAP_SERVERS="localhost:19092,localhost:19094,localhost:19096"
export TOPIC_NAME="load-test-topic"

export PRODUCER_ACKS="1"
export PRODUCER_BATCH_SIZE="32768"
export PRODUCER_LINGER_MS="5"
export MESSAGE_SIZE_BYTES="512"

export CONSUMER_GROUP="load-test-consumer-group"
export CONSUMER_COUNT="3"

echo "Running Kafka load test with k6..."
$K6_EXECUTABLE run --vus 50 --duration 1m "$SCRIPT_NAME"
