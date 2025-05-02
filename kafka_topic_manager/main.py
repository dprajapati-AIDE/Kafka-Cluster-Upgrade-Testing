import os
import sys
import argparse

from logger import get_logger
from manager import KafkaTopicManager

logger = get_logger()


def parse_arguments():
    parser = argparse.ArgumentParser(
        description='Kafka Topic Manager - Create and manage Kafka topics from CSV files')

    # Required arguments
    parser.add_argument('--csv-file', required=True,
                        help='Path to the CSV file with topic configurations')
    parser.add_argument('--kafka-bin', required=True,
                        help='Path to Kafka bin directory with the CLI tools')

    # Kafka connection options
    parser.add_argument('--zk', default='localhost', help='Zookeeper hostname')
    parser.add_argument('--zk-port', default='2181', help='Zookeeper port')

    # Topic configuration defaults
    parser.add_argument('--replication-factor', type=int, default=3,
                        help='Default replication factor for new topics if not specified in CSV')

    # Operation modes
    parser.add_argument('--dry-run', type=bool, default=False,
                        help='Print commands without executing them')

    parser.add_argument('--verbose', '-v', action='store_true',
                        help='Enable verbose logging')

    return parser.parse_args()


def main():
    args = parse_arguments()

    if args.verbose:
        logger.setLevel(logging.DEBUG)
        logger.debug("Verbose logging enabled")

    kafka_topics_sh = os.path.join(args.kafka_bin, 'kafka-topics.sh')
    if not os.path.exists(kafka_topics_sh):
        logger.error(f"Kafka tools not found at: {kafka_topics_sh}")
        logger.info("Please provide the correct path to Kafka bin directory")
        sys.exit(1)

    logger.info(f"Using ZooKeeper connection: {args.zk}:{args.zk_port}")
    logger.info(f"Using Kafka bin directory: {args.kafka_bin}")
    logger.info(f"Reading topics from: {args.csv_file}")

    manager = KafkaTopicManager(args)
    manager.manage_topics()


if __name__ == '__main__':
    import logging
    main()
