from kafka import KafkaConsumer
from kafka.admin import KafkaAdminClient, NewTopic, ConfigResource, ConfigResourceType
from kafka.errors import TopicAlreadyExistsError, NoBrokersAvailable, KafkaConfigurationError, InvalidTopicError, NodeNotReadyError

from csv_reader import ReadTopicFromCSV
from logger import get_logger


class KafkaTopicManager:

    def __init__(self, args):
        self.args = args
        self.dry_run = args.dry_run
        self.logger = get_logger()
        self.csv_reader = ReadTopicFromCSV()

        self.bootstrap_servers = args.bootstrap_server

        self.admin_client = None
        if not self.dry_run:
            try:
                self.admin_client = KafkaAdminClient(
                    bootstrap_servers=self.bootstrap_servers,
                    client_id='kafka_topic_manager'
                )
                self.logger.info(
                    "Successfully connected to Kafka Admin Client.")
            except NoBrokersAvailable as e:
                self.logger.error(
                    f"Could not connect to Kafka brokers at {self.bootstrap_servers}: {e}")
                self.logger.error(
                    "Please ensure Kafka brokers are running and reachable.")
                exit(1)
            except KafkaConfigurationError as e:
                self.logger.error(f"Kafka configuration error: {e}")
                exit(1)
            except Exception as e:
                self.logger.error(
                    f"An unexpected error occurred while connecting to Kafka Admin Client: {e}")
                exit(1)

    def list_topics(self):
        if self.dry_run:
            self.logger.info("[DRY RUN] Would list topics.")
            return []
        try:
            topics = self.admin_client.list_topics()
            self.logger.debug(f"Existing topics: {topics}")
            return list(topics)
        except NoBrokersAvailable:
            self.logger.error(
                "No brokers available to list topics. Is the Kafka cluster running?")
            return []
        except Exception as e:
            self.logger.error(f"Failed to list topics: {e}")
            return []

    def describe_topics(self):
        topics_map = {}
        if self.dry_run:
            self.logger.info("[DRY RUN] Would describe topics.")
            return {}

        try:
            # Get list of topic names
            topic_names = self.list_topics()
            if not topic_names:
                self.logger.info("No topics found in cluster")
                return {}

            self.logger.debug(
                f"Found {len(topic_names)} topics: {topic_names}")

            # Create config resources for all topics
            config_resources = [ConfigResource(ConfigResourceType.TOPIC, topic_name)
                                for topic_name in topic_names]

            # Fetch configurations for all topics
            described_configs = {}
            if config_resources:
                try:
                    config_response = self.admin_client.describe_configs(
                        config_resources)
                    self.logger.debug(
                        f"Config response type: {type(config_response)}")

                    if hasattr(config_response, 'resources') and config_response.resources:
                        for resource_data in config_response.resources:
                            topic_name = resource_data.resource_name
                            configs = {}
                            if hasattr(resource_data, 'config_entries') and resource_data.config_entries:
                                for entry in resource_data.config_entries:
                                    configs[entry.name] = entry.value
                            described_configs[topic_name] = configs

                except Exception as e:
                    self.logger.warning(
                        f"Failed to get topic configurations: {e}")

            consumer = None
            try:
                consumer = KafkaConsumer(
                    bootstrap_servers=self.bootstrap_servers,
                    client_id='kafka_topic_inspector',
                    consumer_timeout_ms=5000
                )

                # Get partition info for each topic
                for topic_name in topic_names:
                    try:
                        # Get partition metadata for this topic
                        partitions_metadata = consumer.partitions_for_topic(
                            topic_name)
                        if partitions_metadata:
                            partitions = len(partitions_metadata)
                        else:
                            partitions = 1

                        replication_factor = 1
                        try:
                            cluster_metadata = consumer._client.cluster
                            if hasattr(cluster_metadata, 'topics') and topic_name in cluster_metadata.topics:
                                topic_metadata = cluster_metadata.topics[topic_name]
                                if topic_metadata.partitions:
                                    first_partition = list(
                                        topic_metadata.partitions.values())[0]
                                    replication_factor = len(
                                        first_partition.replicas)
                        except Exception as e:
                            self.logger.debug(
                                f"Could not get replication factor for {topic_name}: {e}")

                        configs = described_configs.get(topic_name, {})

                        topics_map[topic_name] = {
                            "topic": topic_name,
                            "partitions": partitions,
                            "replication_factor": replication_factor,
                            "configs": configs
                        }

                    except Exception as e:
                        self.logger.warning(
                            f"Error processing topic {topic_name}: {e}")
                        topics_map[topic_name] = {
                            "topic": topic_name,
                            "partitions": 1,
                            "replication_factor": 1,
                            "configs": described_configs.get(topic_name, {})
                        }

            except Exception as e:
                self.logger.warning(
                    f"Could not create consumer for metadata: {e}")
                for topic_name in topic_names:
                    topics_map[topic_name] = {
                        "topic": topic_name,
                        "partitions": 1,
                        "replication_factor": 1,
                        "configs": described_configs.get(topic_name, {})
                    }
            finally:
                if consumer:
                    consumer.close()

            return topics_map

        except NoBrokersAvailable:
            self.logger.error(
                "No brokers available to describe topics. Is the Kafka cluster running?")
            return {}
        except Exception as e:
            self.logger.error(f"Failed to describe topics: {e}")
            return {}

    def create_topic(self, topic, partitions, replication_factor, configs):
        new_topic = NewTopic(
            name=topic,
            num_partitions=partitions,
            replication_factor=replication_factor,
            topic_configs=configs
        )
        if self.dry_run:
            self.logger.info(f"[DRY RUN] Would create topic {topic} with partitions={partitions}, "
                             f"replication_factor={replication_factor}, configs={configs}")
            return True

        try:
            self.admin_client.create_topics([new_topic], validate_only=False)
            self.logger.info(f"Successfully created topic {topic}")
            return True
        except TopicAlreadyExistsError:
            self.logger.warning(
                f"Topic {topic} already exists. Skipping creation.")
            return True
        except NoBrokersAvailable:
            self.logger.error(
                f"No brokers available to create topic {topic}. Is the Kafka cluster running?")
            return False
        except Exception as e:
            self.logger.error(f"Failed to create topic {topic}: {e}")
            return False

    def update_topic_partitions(self, topic, partitions):
        if self.dry_run:
            self.logger.info(
                f"[DRY RUN] Would update partitions for topic {topic} to {partitions}.")
            return True

        try:
            self.logger.info(
                f"Attempting to increase partitions for topic {topic} to {partitions}")
            new_partitions_spec = {topic: partitions}
            self.admin_client.create_partitions(new_partitions_spec)

            self.logger.info(
                f"Successfully updated partitions for topic {topic} to {partitions}")
            return True
        except InvalidTopicError:
            self.logger.error(
                f"Topic {topic} does not exist. Cannot update partitions.")
            return False
        except NoBrokersAvailable:
            self.logger.error(
                f"No brokers available to update partitions for topic {topic}. Is the Kafka cluster running?")
            return False
        except Exception as e:
            self.logger.error(
                f"Failed to update partitions for topic {topic}: {e}")
            return False

    def update_topic_config(self, topic, configs):
        if not configs:
            self.logger.debug(f"No configs to update for topic {topic}.")
            return True

        if self.dry_run:
            self.logger.info(
                f"[DRY RUN] Would update config for topic {topic}: {configs}")
            return True

        try:
            config_resources = [ConfigResource(
                ConfigResourceType.TOPIC, topic, configs)]
            self.admin_client.alter_configs(config_resources)
            self.logger.info(
                f"Successfully updated config for topic {topic}: {configs}")
            return True
        except InvalidTopicError:
            self.logger.error(
                f"Topic {topic} does not exist. Cannot update configs.")
            return False
        except NoBrokersAvailable:
            self.logger.error(
                f"No brokers available to update config for topic {topic}. Is the Kafka cluster running?")
            return False
        except Exception as e:
            self.logger.error(
                f"Failed to update config for topic {topic}: {e}")
            return False

    def delete_topic(self, topic):
        if self.dry_run:
            self.logger.info(f"[DRY RUN] Would delete topic {topic}.")
            return True

        try:
            self.admin_client.delete_topics([topic])
            self.logger.info(f"Successfully deleted topic {topic}")
            return True
        except InvalidTopicError:
            self.logger.warning(
                f"Topic {topic} does not exist. Skipping deletion.")
            return True
        except NoBrokersAvailable:
            self.logger.error(
                f"No brokers available to delete topic {topic}. Is the Kafka cluster running?")
            return False
        except Exception as e:
            self.logger.error(f"Failed to delete topic {topic}: {e}")
            return False

    def manage_topics(self):
        default_replication_factor = self.args.replication_factor

        # describe topics
        existing_topics = self.describe_topics()
        self.logger.debug(f"Existing topics: {list(existing_topics.keys())}")

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
                        f"Topic {topic_name} has {existing_info['partitions']} partitions, cannot decrease to {topic_info['partitions']}. "
                        "Kafka does not support decreasing partitions."
                    )

                # Check if configs need to be updated
                current_configs = existing_info['configs']
                new_configs = topic_info['configs']

                configs_to_apply = {
                    k: v for k, v in new_configs.items() if v is not None and v != ''}

                if configs_to_apply != current_configs:
                    filtered_current_configs = {
                        k: current_configs[k] for k in configs_to_apply if k in current_configs}

                    if configs_to_apply != filtered_current_configs:
                        self.logger.info(
                            f"Updating configs for topic {topic_name}")
                        self.logger.debug(f"Current: {current_configs}")
                        self.logger.debug(f"New: {new_configs}")
                        self.logger.debug(
                            f"Configs to apply: {configs_to_apply}")
                        self.update_topic_config(topic_name, configs_to_apply)
                    else:
                        self.logger.info(
                            f"No relevant config changes for topic {topic_name}.")
                else:
                    self.logger.info(
                        f"No config changes needed for topic {topic_name}.")
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
