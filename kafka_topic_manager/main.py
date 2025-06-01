import os
import sys
import argparse
import logging

from logger import get_logger
from manager import KafkaTopicManager

logger = get_logger()


def parse_arguments():
    parser = argparse.ArgumentParser(
        description='Kafka Topic Manager - Create and manage Kafka topics from CSV files using kafka-python')

    # Required arguments
    parser.add_argument('--csv-file', required=True,
                        help='Path to the CSV file with topic configurations')

    # Kafka connection options
    parser.add_argument('--bootstrap-server', required=True,
                        help='Kafka bootstrap server address (e.g., localhost:9092)')
    # Zookeeper argument is now deprecated for KafkaAdminClient, but we keep it for backward compatibility
    # and log a warning if used.
    parser.add_argument('--zookeeper', default=None,
                        help='Deprecated: Zookeeper hostname. KafkaAdminClient uses --bootstrap-server.')

    # Topic configuration defaults
    parser.add_argument('--replication-factor', type=int, default=3,
                        help='Default replication factor for new topics if not specified in CSV')

    # Operation modes
    parser.add_argument('--dry-run', action='store_true',
                        help='Print commands without executing them')

    parser.add_argument('--verbose', '-v', action='store_true',
                        help='Enable verbose logging')

    return parser.parse_args()


def main():
    args = parse_arguments()

    if args.verbose:
        logger.setLevel(logging.DEBUG)
        logger.debug("Verbose logging enabled")

    if not args.bootstrap_server:
        logger.error(
            "Missing --bootstrap-server argument. This is required for kafka-python.")
        sys.exit(1)

    if args.zookeeper:
        logger.warning("The --zookeeper argument is deprecated when using kafka-python for topic management. "
                       "Please use --bootstrap-server instead.")

    logger.info(f"Using Kafka bootstrap server: {args.bootstrap_server}")
    logger.info(f"Reading topics from: {args.csv_file}")

    manager = KafkaTopicManager(args)
    manager.manage_topics()


if __name__ == '__main__':
    main()
