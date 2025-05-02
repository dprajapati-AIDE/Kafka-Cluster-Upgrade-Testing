import re

from csv_reader import ReadTopicFromCSV
from logger import get_logger
from utils import execute_command, get_broker_connection

# Regular expression to extract topic information from kafka-topics --describe output
RE_TOPIC_DESCRIPTION = re.compile(
    r"Topic:(?P<topic>\S+)\s+PartitionCount:(?P<partitions>\d+)\s+ReplicationFactor:(?P<replication_factor>\d+)\s+Configs:(?P<configs>.+)?")


class KafkaTopicManager:

    def __init__(self, args):
        self.args = args
        self.kafka_bin = args.kafka_bin
        self.dry_run = args.dry_run
        self.logger = get_logger()
        self.broker = get_broker_connection(args.zk, args.zk_port)
        self.csv_reader = ReadTopicFromCSV()

    def list_topics(self):
        cmd = [f"{self.kafka_bin}/kafka-topics.sh",
               "--zookeeper", self.broker, "--list"]

        exit_code, output = execute_command(cmd, self.args.dry_run)

        if exit_code != 0:
            self.logger.error(f"Failed to list topics: {output}")
            return []

        topics = [topic for topic in output.strip().split('\n') if topic]
        self.logger.debug(f"Existing topics: {topics}")
        return topics

    def describe_topics(self):
        topics_map = {}

        cmd = [f"{self.kafka_bin}/kafka-topics.sh",
               "--zookeeper", self.broker, "--describe"]

        exit_code, output = execute_command(cmd, self.args.dry_run)

        if exit_code != 0:
            self.logger.error(f"Failed to describe topics: {output}")
            return {}

        for line in output.strip().split('\n'):
            if not line:
                continue

            result = RE_TOPIC_DESCRIPTION.search(line)
            if result:
                topic = result.group("topic").strip()
                configs_str = result.group("configs") or ""

                topics_map[topic] = {
                    "topic": topic,
                    "partitions": int(result.group("partitions")),
                    "replication_factor": int(result.group("replication_factor")),
                    "configs": self.csv_reader.parse_topic_configs(configs_str)
                }

        return topics_map

    def create_topic(self, topic, partitions, replication_factor, configs):
        cmd = [
            f"{self.kafka_bin}/kafka-topics.sh",
            "--zookeeper", self.broker,
            "--create",
            "--topic", topic,
            "--partitions", str(partitions),
            "--replication-factor", str(replication_factor)
        ]

        self.logger.info(
            f"Creating topic: {topic} with partitions: {partitions}, replication-factor: {replication_factor}")
        exit_code, output = execute_command(cmd, self.args.dry_run)

        if exit_code != 0:
            self.logger.error(f"Failed to create topic {topic}: {output}")
            return False

        if configs and len(configs) > 0:
            return self.update_topic_config(topic, configs)

        return True

    def update_topic_partitions(self, topic, partitions):
        cmd = [
            f"{self.kafka_bin}/kafka-topics.sh",
            "--zookeeper", self.broker,
            "--alter",
            "--topic", topic,
            "--partitions", str(partitions)
        ]

        self.logger.info(
            f"Updating partitions for topic {topic} to {partitions}")
        exit_code, output = execute_command(cmd, self.args.dry_run)

        if exit_code != 0:
            self.logger.error(
                f"Failed to update partitions for topic {topic}: {output}")
            return False

        return True

    def format_configs_string(self, configs):
        return ",".join(f"{k}={v}" for k, v in sorted(configs.items()))

    def update_topic_config(self, topic, configs):
        if not configs:
            return True

        config_str = self.format_configs_string(configs)

        kafka_configs_sh = f"{self.kafka_bin}/kafka-configs.sh"

        try:
            import os
            if os.path.exists(kafka_configs_sh):
                cmd = [
                    kafka_configs_sh,
                    "--zookeeper", self.broker,
                    "--entity-type", "topics",
                    "--entity-name", topic,
                    "--alter",
                    "--add-config", config_str
                ]
            else:
                self.logger.warning(
                    f"kafka-configs.sh not found, falling back to kafka-topics.sh for config updates")
                cmd = [
                    f"{self.kafka_bin}/kafka-topics.sh",
                    "--zookeeper", self.broker,
                    "--alter",
                    "--topic", topic,
                    "--config", config_str
                ]
        except Exception as e:
            self.logger.warning(
                f"Error checking for kafka-configs.sh: {e}, using kafka-topics.sh instead")
            cmd = [
                f"{self.kafka_bin}/kafka-topics.sh",
                "--zookeeper", self.broker,
                "--alter",
                "--topic", topic,
                "--config", config_str
            ]

        self.logger.info(f"Updating config for topic {topic}: {config_str}")
        exit_code, output = execute_command(cmd, self.args.dry_run)

        if exit_code != 0:
            self.logger.error(
                f"Failed to update config for topic {topic}: {output}")
            return False

        return True

    def delete_topic(self, topic):
        cmd = [
            f"{self.kafka_bin}/kafka-topics.sh",
            "--zookeeper", self.broker,
            "--delete",
            "--topic", topic
        ]
        exit_code, output = execute_command(cmd, self.args.dry_run)

        if exit_code != 0:
            self.logger.error(f"Failed to delete topic {topic}: {output}")
            return False

        return True

    def manage_topics(self):
        default_replication_factor = self.args.replication_factor
        if not self.args.csv_file:
            self.logger.error("CSV file not specified")
            return

        # describe topics
        try:
            existing_topics = self.describe_topics()
            self.logger.debug(
                f"Existing topics: {list(existing_topics.keys())}")
        except Exception as e:
            self.logger.error(f"Error getting existing topics: {e}")
            existing_topics = {}

        try:
            topics_to_manage = self.csv_reader.read_csv_topics(
                self.args.csv_file, default_replication_factor)
            if topics_to_manage:
                self.logger.info(
                    f"Found {len(topics_to_manage)} topics in CSV file")
            else:
                self.logger.error(
                    "No topics found in CSV file or error reading file")
                return
        except Exception as e:
            self.logger.error(f"Error reading CSV file: {e}")
            return

        # process each topic
        for topic_info in topics_to_manage:
            topic_name = topic_info['topic']

            if topic_name in existing_topics:
                existing_info = existing_topics[topic_name]

                # Check if partitions need to be updated
                if topic_info['partitions'] > existing_info['partitions']:
                    self.logger.info(
                        f"Updating partitions for topic {topic_name} from {existing_info['partitions']} to {topic_info['partitions']}")
                    self.update_topic_partitions(
                        topic_name, topic_info['partitions'])
                elif topic_info['partitions'] < existing_info['partitions']:
                    self.logger.warning(
                        f"Topic {topic_name} has {existing_info['partitions']} partitions, cannot decrease to {topic_info['partitions']}")

                # Check if configs need to be updated
                current_configs = existing_info['configs']
                new_configs = topic_info['configs']

                if current_configs != new_configs:
                    self.logger.info(
                        f"Updating configs for topic {topic_name}")
                    self.logger.debug(f"Current: {current_configs}")
                    self.logger.debug(f"New: {new_configs}")
                    self.update_topic_config(topic_name, new_configs)
            else:
                # Create new topic
                self.logger.info(
                    f"Creating topic {topic_name} with {topic_info['partitions']} partitions and replication factor {topic_info['replication_factor']}")
                self.create_topic(
                    topic_name,
                    topic_info['partitions'],
                    topic_info['replication_factor'],
                    topic_info['configs']
                )
