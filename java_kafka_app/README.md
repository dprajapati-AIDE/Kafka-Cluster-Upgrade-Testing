# Java Kafka Producer-Consumer

A modular, CLI-driven Kafka producer-consumer application written in Java Maven. It uses CLI commands and supports YAML-based configuration, structured logging, and clean package organization.

## Setting up Multiple Java Version

### Installation and Configuration of jenv

Install `jenv`
```
brew install jenv
```

Add the following lines to your shell startup file
```
export PATH="$HOME/.jenv/bin:$PATH"
eval "$(jenv init -)"
```

Then, source the file to apply the changes:
```
source ~/.zshrc
```


### Installation of Java Versions

Install Java8 using below command
```
brew install --cask temurin@8
```

Install Java11 using below command
```
brew install openjdk@11
```


### Configuration

Add Java8 to jenv
```
jenv add /Library/Java/JavaVirtualMachines/temurin-8.jdk/Contents/Home
```

Setup symlink for Java11 and add to jenv
```
sudo ln -sfn /opt/homebrew/opt/openjdk@11/libexec/openjdk.jdk /Library/Java/JavaVirtualMachines/openjdk-11.jdk
```

```
jenv add /opt/homebrew/Cellar/openjdk@11/11.0.27/libexec/openjdk.jdk/Contents/Home
```

#### List all the versions
```
jenv versions
```

#### List current version
```
jenv version
```

### Set Global and Local versions
The global version will be the default for your system, while the local version will be specific to the current directory.

```
jenv global 11.0
```

```
jenv local 11.0
```

## Create Maven Project

Update details like groupId and artifactId based on requirement
```
mvn archetype:generate -DgroupId=com.kafka.app -DartifactId=java_kafka_app -DarchetypeArtifactId=maven-archetype-quickstart -DinteractiveMode=false
```

### To compile for Java 8 (default):
```
jenv local 1.8
```
```
mvn clean install
```

### To compile for Java 11:
```
jenv local 11
```
```
mvn clean install -Djava.version=11
```

### To compile for Java 21:
```
jenv local 21
```
```
mvn clean install -Djava.version=21
```


## Project Structure

```
.
├── README.md
├── dependency-reduced-pom.xml                     # Maven-generated POM (for shaded jars)
├── pom.xml
├── src/
│   └── main/
│       ├── java/
│       │   └── com/kafka/app/
│       │       ├── App.java                       # Legacy entry point (can be deprecated)
│       │       ├── Run.java                       # New CLI-based main class with role, msg-count, consumer-group support
│       │       ├── consumer/                      # Kafka Consumer components
│       │       │   ├── Consumer.java              # Handles message processing and deserialization
│       │       │   └── ConsumerGroup.java         # Kafka consumer group setup and polling loop
│       │       ├── kafka/
│       │       │   └── KafkaManager.java          # Checks Kafka broker connectivity
│       │       ├── logger/
│       │       │   └── AppLogger.java             # SLF4J logger configuration
│       │       ├── model/                         # Configuration and domain models
│       │       │   ├── AppConfig.java             # Root config mapping for YAML
│       │       │   ├── DeviceConfig.java          # Vendor + device types config
│       │       │   ├── DeviceMessageModel.java    # Kafka message payload structure
│       │       │   ├── DeviceTypes.java           # Device types (firewall, router, switch)
│       │       │   ├── KafkaCluster.java          # Kafka cluster config (brokers, topics, version)
│       │       │   ├── KafkaConfig.java           # Wrapper for multiple clusters
│       │       │   ├── LoggingConfig.java         # Logging configuration
│       │       │   └── Topic.java                 # Kafka topic config model
│       │       ├── producer/
│       │       │   └── Producer.java              # Kafka producer logic for all clusters/topics
│       │       └── topic/
│       │           └── TopicManager.java          # Ensures required topics exist
│       └── resources/
│           ├── config.yaml                        # Application config
│           └── log4j2.properties                  # Log4j logging configuration

```

## Running the Application

### Build the jar file
```
mvn clean install
```

### Producer Mode
```
java -jar target/java_producer_consumer-shaded.jar --role producer --msg-count 10
```

`--msg-count` is optional; defaults to the configured value if not provided.

### Consumer Mode
```
java -jar target/java_producer_consumer-shaded.jar --role consumer --consumer-group demo
```

## Producer Consumer Together Mode
```
java -jar target/java_producer_consumer-shaded.jar --role both --msg-count 10 --consumer-group test-both
```

## Configuration
Edit the `resources/config.yaml` file to define Kafka settings, topic names and devices related configurations.

## Project Dependencies
| Dependency                                                 | Version  | Purpose                                                        |
| ---------------------------------------------------------- | -------- | -------------------------------------------------------------- |
| `junit:junit`                                              | 4.11     | Unit testing framework (used for writing test cases)           |
| `com.fasterxml.jackson.core:jackson-databind`              | 2.17.0   | Core Jackson library for JSON/YAML parsing                     |
| `com.fasterxml.jackson.dataformat:jackson-dataformat-yaml` | 2.17.0   | Enables reading `YAML` configuration files                     |
| `org.slf4j:slf4j-api`                                      | 1.7.36   | Logging facade used throughout the app                         |
| `org.apache.logging.log4j:log4j-api`                       | 2.20.0   | Core Log4j API for structured logging                          |
| `org.apache.logging.log4j:log4j-core`                      | 2.20.0   | Logging backend implementation                                 |
| `org.apache.logging.log4j:log4j-slf4j-impl`                | 2.20.0   | Bridges SLF4J logging to Log4j2 backend                        |
| `org.apache.kafka:kafka-clients`                           | 0.11.0.3 | Kafka client library for Java (consumer, producer, admin APIs) |
| `com.github.javafaker:javafaker`                           | 1.0.2    | Random data generator used for simulating device messages      |

## Build Plugins
| Plugin                  | Version | Purpose                                                             |
| ----------------------- | ------- | ------------------------------------------------------------------- |
| `maven-compiler-plugin` | 3.8.1   | Compiles Java source code using specified Java version (`1.8`)      |
| `maven-jar-plugin`      | 3.3.0   | Builds standard JAR with custom manifest entries                    |
| `maven-shade-plugin`    | 3.5.2   | Creates an executable “fat” JAR including all dependencies (shaded) |


---

#### If you are using VS Code and get error regarding Class Path
- Go to settings or open it using `command + ,`
- Search for `java.project.sourcePaths`
- Add `"src/main/java"`
- Reload the workspace.
