#!/bin/bash

# Environment variables with defaults
ZOOKEEPER_SERVER_ID=${ZOOKEEPER_SERVER_ID:-1}
ZOOKEEPER_TICK_TIME=${ZOOKEEPER_TICK_TIME:-2000}
ZOOKEEPER_INIT_LIMIT=${ZOOKEEPER_INIT_LIMIT:-5}
ZOOKEEPER_SYNC_LIMIT=${ZOOKEEPER_SYNC_LIMIT:-2}
ZOOKEEPER_SERVERS=${ZOOKEEPER_SERVERS:-"zookeeper1:2888:3888,zookeeper2:2888:3888,zookeeper3:2888:3888"}

# Ensure necessary directories exist
mkdir -p /opt/zookeeper/data

# Log the configuration
echo "Starting ZooKeeper with the following configuration:"
echo "  ZOOKEEPER_SERVER_ID=${ZOOKEEPER_SERVER_ID}"
echo "  ZOOKEEPER_TICK_TIME=${ZOOKEEPER_TICK_TIME}"
echo "  ZOOKEEPER_INIT_LIMIT=${ZOOKEEPER_INIT_LIMIT}"
echo "  ZOOKEEPER_SYNC_LIMIT=${ZOOKEEPER_SYNC_LIMIT}"
echo "  ZOOKEEPER_SERVERS=${ZOOKEEPER_SERVERS}"

# Create dynamic Zookeeper configuration
cat <<EOF > /opt/zookeeper/conf/zoo.cfg
tickTime=${ZOOKEEPER_TICK_TIME}
dataDir=/opt/zookeeper/data
clientPort=2181
initLimit=${ZOOKEEPER_INIT_LIMIT}
syncLimit=${ZOOKEEPER_SYNC_LIMIT}
EOF

# Set up the servers in zoo.cfg
IFS=',' read -ra SERVER_ARRAY <<< "$ZOOKEEPER_SERVERS"
for server in "${SERVER_ARRAY[@]}"; do
    ID=$(echo $server | cut -d':' -f1 | sed 's/zookeeper//')
    echo "server.$ID=$server" >> /opt/zookeeper/conf/zoo.cfg
done

# Set Zookeeper ID for the current node
echo "${ZOOKEEPER_SERVER_ID}" > /opt/zookeeper/data/myid
echo "ZooKeeper ID set to ${ZOOKEEPER_SERVER_ID}"

# Start Zookeeper server
echo "Starting ZooKeeper..."
exec "$@"
