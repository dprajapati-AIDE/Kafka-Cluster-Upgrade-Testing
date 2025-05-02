#!/bin/bash

KAFKA_BROKER_ID=${KAFKA_BROKER_ID:-1}
KAFKA_LISTENERS=${KAFKA_LISTENERS:-"INTERNAL://0.0.0.0:19091,EXTERNAL://0.0.0.0:19092"}
KAFKA_ADVERTISED_LISTENERS=${KAFKA_ADVERTISED_LISTENERS:-"INTERNAL://kafka-1-372:19091,EXTERNAL://localhost:19092"}
KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=${KAFKA_LISTENER_SECURITY_PROTOCOL_MAP:-"INTERNAL:PLAINTEXT,EXTERNAL:PLAINTEXT"}
KAFKA_INTER_BROKER_LISTENER_NAME=${KAFKA_INTER_BROKER_LISTENER_NAME:-"INTERNAL"}
KAFKA_ZOOKEEPER_CONNECT=${KAFKA_ZOOKEEPER_CONNECT:-"zookeeper-1-372:2181,zookeeper-2-372:2181,zookeeper-3-372:2181"}
KAFKA_LOG_DIRS=${KAFKA_LOG_DIRS:-"/var/lib/kafka/data"}

echo "Starting Kafka 3.7.2 with the following configuration:"
echo "  KAFKA_BROKER_ID=${KAFKA_BROKER_ID}"
echo "  KAFKA_LISTENERS=${KAFKA_LISTENERS}"
echo "  KAFKA_ADVERTISED_LISTENERS=${KAFKA_ADVERTISED_LISTENERS}"
echo "  KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=${KAFKA_LISTENER_SECURITY_PROTOCOL_MAP}"
echo "  KAFKA_INTER_BROKER_LISTENER_NAME=${KAFKA_INTER_BROKER_LISTENER_NAME}"
echo "  KAFKA_ZOOKEEPER_CONNECT=${KAFKA_ZOOKEEPER_CONNECT}"
echo "  KAFKA_LOG_DIRS=${KAFKA_LOG_DIRS}"

mkdir -p ${KAFKA_LOG_DIRS}

cat <<EOF > /opt/kafka/config/server.properties
# Broker settings
broker.id=${KAFKA_BROKER_ID}
listeners=${KAFKA_LISTENERS}
advertised.listeners=${KAFKA_ADVERTISED_LISTENERS}
listener.security.protocol.map=${KAFKA_LISTENER_SECURITY_PROTOCOL_MAP}
inter.broker.listener.name=${KAFKA_INTER_BROKER_LISTENER_NAME}
zookeeper.connect=${KAFKA_ZOOKEEPER_CONNECT}
log.dirs=${KAFKA_LOG_DIRS}

# Topic defaults
num.partitions=3
default.replication.factor=2
min.insync.replicas=2
auto.create.topics.enable=true

# Performance tuning
num.network.threads=3
num.io.threads=8
socket.send.buffer.bytes=102400
socket.receive.buffer.bytes=102400
socket.request.max.bytes=104857600
group.initial.rebalance.delay.ms=0

# Log settings
log.retention.hours=168
log.segment.bytes=1073741824
log.retention.check.interval.ms=300000
EOF

echo "Starting Kafka..."
exec "$@"
