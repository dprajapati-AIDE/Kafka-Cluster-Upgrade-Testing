# Kafka Cluster Upgrade Testing

This repository provides resources for setting up and testing Kafka clusters with versions 2.1.1 and 3.7.2 using custom Docker images. It also includes a Kafka UI for easy cluster management and monitoring.

---

### 📦 Project Structure

```
.
├── README.md
├── kafka-ui/
│   └── docker-compose.yaml        # Kafka UI setup
├── kafka_2.1.1/                   # Kafka 2.1.1 cluster setup
│   ├── Dockerfile.kafka
│   ├── Dockerfile.zookeeper
│   ├── docker-compose.yaml
│   ├── entrypoint-kafka.sh
│   └── entrypoint-zookeeper.sh
└── kafka_3.7.2/                   # Kafka 3.7.2 cluster setup
    ├── Dockerfile.kafka
    ├── Dockerfile.zookeeper
    ├── docker-compose.yaml
    ├── entrypoint-kafka.sh
    └── entrypoint-zookeeper.sh
```

### 🛠️ Build Docker Images

To build Kafka and Zookeeper Docker images for each version:

- 1. Navigate to the version-specific directory:
    ```
    cd kafka_<version>
    ```
- 2. Build Kafka image:
    ```
    docker build -t kafka:<version> -f Dockerfile.kafka .
    ```
- 3. Build Zookeeper image:
    ```
    docker build -t zookeeper:<version> -f Dockerfile.zookeeper .
    ```
    > Replace <version> with either 2.1.1 or 3.7.2 as needed.

### 🚀 Run Kafka Clusters

To start the Kafka clusters:

- 1. From each version directory, run:
    ```
    docker compose up -d
    ```
- 2. Once both clusters are running, start the Kafka UI:
    ```
    cd kafka-ui
    docker compose up -d
    ```

### 📊 Kafka UI
The Kafka UI provides a web interface to manage and monitor Kafka clusters. It connects to both 2.1.1 and 3.7.2 clusters by default as configured in the kafka-ui/docker-compose.yaml file.

### 🚧 Next Steps

1. Create Topics using CSV File
You can automate the process of creating Kafka topics across both Kafka clusters (version 2.1.1 and 3.7.2) using a Python script. The script will:

- Accept a CSV file as input that defines topic configurations such as topic name, number of partitions, replication factor, retention policy, and cleanup policy.

- Use the `argparse` module to accept cluster DNS or IP addresses as arguments for the Kafka clusters.

- Use Python's Kafka admin libraries (like confluent_kafka or kafka-python) to programmatically create topics on both clusters.

- Example CSV Format
    ```
    topic_name,num_partitions,replication_factor,retention_ms,cleanup_policy
    test-topic-1,3,1,60000,delete
    test-topic-2,5,1,120000,compact
    ```

2. Kafka Producer-Consumer Applications (Python, Java, Go)
Once topics are created, it's time to generate concurrent load on both Kafka clusters using producer-consumer applications in Python, Java, and Go. Each application will:

- Generate random data and produce it to the topics created in Step 1.

- Consume data from the same topics to simulate load.

- Run concurrently for both Kafka clusters (version 2.1.1 and 3.7.2).

    Each language-specific implementation will use environment variables or configuration files to determine the cluster DNS/IP addresses and other necessary configurations.

    Key Tasks:

    1. Python Producer-Consumer:

    - Use the confluent_kafka or kafka-python libraries.

    - Use Python's concurrent.futures or asyncio for concurrent production and consumption.

    2. Java Producer-Consumer:

    - Use KafkaProducer and KafkaConsumer from the Kafka Java client.

    - Use Java's concurrency mechanisms (e.g., ExecutorService) for parallel processing.

    3. Go Producer-Consumer:

    - Use the github.com/segmentio/kafka-go library.

    - Use Go's concurrency model (goroutines) for parallel execution.

    Each producer-consumer app will run in parallel for both Kafka clusters, simulating real-time load and testing system performance.
