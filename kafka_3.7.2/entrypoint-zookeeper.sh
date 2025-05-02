#!/bin/bash

set -e

ZOOKEEPER_SERVER_ID=${ZOOKEEPER_SERVER_ID:-1}
ZOOKEEPER_TICK_TIME=${ZOOKEEPER_TICK_TIME:-2000}
ZOOKEEPER_INIT_LIMIT=${ZOOKEEPER_INIT_LIMIT:-5}
ZOOKEEPER_SYNC_LIMIT=${ZOOKEEPER_SYNC_LIMIT:-2}
ZOOKEEPER_4LW_COMMANDS_WHITELIST=${ZOOKEEPER_4LW_COMMANDS_WHITELIST:-"ruok,stat,mntr"}
ZOOKEEPER_SERVERS=${ZOOKEEPER_SERVERS:-"zookeeper-1-372:2888:3888,zookeeper-2-372:2888:3888,zookeeper-3-372:2888:3888"}

mkdir -p /opt/zookeeper/data

echo "Starting ZooKeeper 3.8.3 with the following configuration:"
echo "  ZOOKEEPER_SERVER_ID=${ZOOKEEPER_SERVER_ID}"
echo "  ZOOKEEPER_TICK_TIME=${ZOOKEEPER_TICK_TIME}"
echo "  ZOOKEEPER_INIT_LIMIT=${ZOOKEEPER_INIT_LIMIT}"
echo "  ZOOKEEPER_SYNC_LIMIT=${ZOOKEEPER_SYNC_LIMIT}"
echo "  ZOOKEEPER_SERVERS=${ZOOKEEPER_SERVERS}"
echo "  ZOOKEEPER_4LW_COMMANDS_WHITELIST=${ZOOKEEPER_4LW_COMMANDS_WHITELIST}"

cat <<EOF > /opt/zookeeper/conf/zoo.cfg
# Basic configurations
tickTime=${ZOOKEEPER_TICK_TIME}
dataDir=/opt/zookeeper/data
clientPort=2181
initLimit=${ZOOKEEPER_INIT_LIMIT}
syncLimit=${ZOOKEEPER_SYNC_LIMIT}
4lw.commands.whitelis=${ZOOKEEPER_4LW_COMMANDS_WHITELIST}

# Performance tuning
snapCount=100000
autopurge.snapRetainCount=3
autopurge.purgeInterval=1
maxClientCnxns=60
admin.enableServer=true

# Metrics configuration
metricsProvider.className=org.apache.zookeeper.metrics.prometheus.PrometheusMetricsProvider
metricsProvider.httpPort=7000
EOF

IFS=',' read -ra SERVER_ARRAY <<< "$ZOOKEEPER_SERVERS"
ID=1
for server in "${SERVER_ARRAY[@]}"; do
    echo "server.${ID}=${server}" >> /opt/zookeeper/conf/zoo.cfg
    ((ID++))
done

echo "${ZOOKEEPER_SERVER_ID}" > /opt/zookeeper/data/myid
echo "ZooKeeper ID set to ${ZOOKEEPER_SERVER_ID}"

echo "Starting ZooKeeper..."
exec "$@"
