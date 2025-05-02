# Kafka Cluster Upgrade Testing

This repository provides resources for setting up and testing Kafka clusters with versions 2.1.1 and 3.7.2 using custom Docker images. It also includes a Kafka UI for easy cluster management and monitoring.

---

### 📦 Project Structure

```
.
├── README.md
├── kafka_2.1.1                             # Kafka 2.1.1 configuration
│   ├── Dockerfile.kafka
│   ├── Dockerfile.zookeeper
│   ├── docker-compose.yaml
│   ├── entrypoint-kafka.sh
│   └── entrypoint-zookeeper.sh
├── kafka_3.7.2                             # Kafka 3.7.2 configuration
│   ├── Dockerfile.kafka
│   ├── Dockerfile.zookeeper
│   ├── docker-compose.yaml
│   ├── entrypoint-kafka.sh
│   └── entrypoint-zookeeper.sh
├── kafka_topic_manager                     # Automation to create multiple topics on cluster
│   ├── __init__.py
│   ├── csv_reader.py
│   ├── logger.py
│   ├── main.py
│   ├── manager.py
│   ├── topics
│   │   └── topics.csv
│   └── utils.py
└── kafka_ui                                # User Interface to manage multiple kafka cluster
    └── docker-compose.yaml
```

### 🛠️ Build Docker Images

- To build Kafka and Zookeeper Docker images for each version:

    1. Navigate to the version-specific directory:
        ```
        cd kafka_<version>
        ```

    2. Build Kafka image:
        ```
        docker build -t kafka:<version> -f Dockerfile.kafka .
        ```

    3. Build Zookeeper image:
        ```
        docker build -t zookeeper:<version> -f Dockerfile.zookeeper .
        ```
    > Replace <version> with either 2.1.1 or 3.7.2 as needed.

### 🚀 Run Kafka Clusters

- To start the Kafka clusters:

    1. From each version directory, run:
        ```
        docker compose up -d
        ```

    2. Once both clusters are running, start the Kafka UI:
        ```
        cd kafka-ui
        docker compose up -d
        ```

### 📊 Kafka UI
- The Kafka UI provides a web interface to manage and monitor Kafka clusters. It connects to both 2.1.1 and 3.7.2 clusters by default as configured in the kafka-ui/docker-compose.yaml file.

### create topics

- #### Docker-Based Setup

    - If you're using a Docker-based Kafka setup, mount the `kafka_topic_manager` module as a volume and run the appropriate command based on your Kafka version:

        > Note: The following commands are based on the DNS configuration used in my Docker Compose setup. You may need to update the Zookeeper or Kafka broker addresses (zookeeper-1-211, kafka-1-372, etc.) to match your own environment.

        **For Kafka 2.x:**
        ```
        python3 -u kafka_topic_manager/main.py --csv-file kafka_topic_manager/topics/topics.csv --kafka-bin /opt/kafka/bin --zookeeper zookeeper-1-211:2181 --replication-factor 1 --verbose
        ```

        **For Kafka 3.x:**
        ```
        python3 -u kafka_topic_manager/main.py --csv-file kafka_topic_manager/topics/topics.csv --kafka-bin /opt/kafka/bin --bootstrap-server kafka-1-372:19091 --replication-factor 1 --verbose
        ```

- #### VM-Based Kafka Cluster

    - If you're working with a VM-based Kafka cluster, ensure the Kafka binaries are installed locally. Then, update the --zookeeper or --bootstrap-server parameters accordingly and run the appropriate command based on your Kafka version.

### 🚧 Next Steps

#### Kafka Producer-Consumer Applications (Python, Java, Go)
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
