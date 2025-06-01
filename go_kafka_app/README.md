# Golang Kafka Producer-Consumer

A modular, CLI-driven Kafka producer-consumer application written in Go. It uses the Cobra library for CLI commands and supports YAML-based configuration, structured logging, and clean package organization.

It also supports implementation of Internal Library - Bus3
- [mist_go_utils](github.com/mistsys/mist_go_utils)
- [sarama-consumer](github.com/mistsys/sarama-consumer)


## 🛠 Installation

Ensure you have Go installed (version 1.23.2). Clone the repository and run:

```bash
go mod tidy
```

This will download and install the necessary dependencies listed in `go.mo` and `go.sum`.


## Project Structure

```
├── cmd
│   └── root.go                # Cobra CLI entry point
├── config
│   └── config.yaml            # Application configuration file
├── go.mod
├── go.sum
├── internal
│   ├── cli
│   │   └── run.go             # CLI argument handling
│   ├── config
│   │   └── config.go          # YAML config loader
│   ├── consumer
│   │   ├── consumer.go        # Kafka consumer logic
│   │   └── group.go           # Consumer group logic
│   ├── kafka
│   │   └── client.go          # Kafka client initialization
│   ├── logger
│   │   └── logger.go          # Logger setup
│   ├── producer
│   │   ├── model.go           # Producer message model
│   │   └── producer.go        # Kafka producer logic
│   ├── topic
│   │   └── topic.go           # Topic creation/management
│   └── utils
│       └── utils.go           # Utility functions
└── main
    └── main.go                # Application entry point
```

## Running the Application

### Producer Mode
```
go run main/main.go --role=producer --msg-count=5 --msg-size-kb=100
```
or to use internal library implementation
```
go run main/main.go --role=producer --msg-count=5 --bus3=true --msg-size-kb=100
```
`--msg-count` and `--msg-size-kb` are optional; defaults to the configured value if not provided.

### Consumer Mode
```
go run main/main.go --role=consumer --consumer-group=demo
```
or to use internal library implementation
```
go run main/main.go --role=consumer --consumer-group=demo --bus3=true
```

## Producer Consumer Together Mode
```
go run main/main.go --role=both
```
or to use internal library implementation
```
go run main/main.go --role=both --msg-count=10 --consumer-group=demo --bus3=true --msg-size-kb=100
```


## Configuration
Edit the `config/config.yaml` file to define Kafka settings, topic names and devices related configurations.

## Core Dependencies
| Library                                       | Purpose                                  |
| --------------------------------------------- | ---------------------------------------- |
| [`sarama`](https://github.com/Shopify/sarama) | Kafka client for Go                      |
| [`cobra`](https://github.com/spf13/cobra)     | CLI command parsing                      |
| [`zap`](https://github.com/uber-go/zap)       | High-performance structured logging      |
| [`faker`](https://github.com/go-faker/faker)  | Generate fake data for producer messages |
|[`mist_go_utils`](github.com/mistsys/mist_go_utils)|Mistsys wrapper of sarama|
|[`sarama-consumer`](github.com/mistsys/sarama-consumer)|Mistsys wrapper of sarama consumer|
