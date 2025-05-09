#!/bin/bash

cat > kafka-jmx-config-1.yml << EOF
---
hostPort: kafka-1-372:9999
lowercaseOutputName: true
lowercaseOutputLabelNames: true
ssl: false
rules:
  - pattern: "kafka.server<type=(.+), name=(.+)PerSec\\\w*><>Count"
    name: "kafka_server_\$1_\$2_total"
    type: COUNTER
  - pattern: "kafka.server<type=(.+), name=(.+)><>Value"
    name: "kafka_server_\$1_\$2"
    type: GAUGE
  - pattern: "kafka.network<type=(.+), name=(.+)><>Value"
    name: "kafka_network_\$1_\$2"
    type: GAUGE
  - pattern: "kafka.cluster<type=(.+), name=(.+)><>Value"
    name: "kafka_cluster_\$1_\$2"
    type: GAUGE
  - pattern: "java.lang<type=(.+), name=(.+)><>Value"
    name: "java_lang_\$1_\$2"
    type: GAUGE
EOF

cat > kafka-jmx-config-2.yml << EOF
---
hostPort: kafka-2-372:9998
lowercaseOutputName: true
lowercaseOutputLabelNames: true
ssl: false
rules:
  - pattern: "kafka.server<type=(.+), name=(.+)PerSec\\\w*><>Count"
    name: "kafka_server_\$1_\$2_total"
    type: COUNTER
  - pattern: "kafka.server<type=(.+), name=(.+)><>Value"
    name: "kafka_server_\$1_\$2"
    type: GAUGE
  - pattern: "kafka.network<type=(.+), name=(.+)><>Value"
    name: "kafka_network_\$1_\$2"
    type: GAUGE
  - pattern: "kafka.cluster<type=(.+), name=(.+)><>Value"
    name: "kafka_cluster_\$1_\$2"
    type: GAUGE
  - pattern: "java.lang<type=(.+), name=(.+)><>Value"
    name: "java_lang_\$1_\$2"
    type: GAUGE
EOF

cat > kafka-jmx-config-3.yml << EOF
---
hostPort: kafka-3-372:9997
lowercaseOutputName: true
lowercaseOutputLabelNames: true
ssl: false
rules:
  - pattern: "kafka.server<type=(.+), name=(.+)PerSec\\\w*><>Count"
    name: "kafka_server_\$1_\$2_total"
    type: COUNTER
  - pattern: "kafka.server<type=(.+), name=(.+)><>Value"
    name: "kafka_server_\$1_\$2"
    type: GAUGE
  - pattern: "kafka.network<type=(.+), name=(.+)><>Value"
    name: "kafka_network_\$1_\$2"
    type: GAUGE
  - pattern: "kafka.cluster<type=(.+), name=(.+)><>Value"
    name: "kafka_cluster_\$1_\$2"
    type: GAUGE
  - pattern: "java.lang<type=(.+), name=(.+)><>Value"
    name: "java_lang_\$1_\$2"
    type: GAUGE
EOF

echo "JMX configuration files created successfully!"
