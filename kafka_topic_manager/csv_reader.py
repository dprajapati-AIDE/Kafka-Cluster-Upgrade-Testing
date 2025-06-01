# csv_reader.py
import os
import csv

from logger import get_logger


class ReadTopicFromCSV:

    def __init__(self):
        self.logger = get_logger()

    def parse_topic_configs(self, configs_str):
        configs = {}
        if not configs_str:
            return configs
        for item in configs_str.split(','):
            item = item.strip()
            if '=' in item:
                key, value = item.split('=', 1)
                configs[key.strip()] = value.strip()
            elif item:
                self.logger.warning(
                    f"Skipping malformed config item: '{item}'")
        return configs

    def read_csv_topics(self, file, default_replication_factor):
        if not os.path.exists(file):
            self.logger.error(f"File Not Found: {file}")
            return []

        topics = []
        try:
            with open(file, 'r') as f:
                reader = csv.DictReader(f)
                for row in reader:
                    if 'topics' not in row or 'partitions' not in row:
                        self.logger.warning(
                            f"Skipping row due to missing 'topics' or 'partitions' field: {row}")
                        continue

                    try:
                        partitions = int(row['partitions'])
                    except ValueError:
                        self.logger.error(
                            f"Invalid partitions value for topic '{row.get('topics', 'N/A')}': {row['partitions']}. Skipping.")
                        continue

                    replication_factor = default_replication_factor
                    if 'replication_factor' in row and row['replication_factor'].strip():
                        try:
                            replication_factor = int(row['replication_factor'])
                        except ValueError:
                            self.logger.warning(
                                f"Invalid replication_factor for topic '{row['topics']}': {row['replication_factor']}. Using default: {default_replication_factor}")

                    configs_str = row.get('configs', '').strip().strip('"')
                    configs = {}
                    if configs_str:
                        configs = self.parse_topic_configs(configs_str)

                    topic_info = {
                        'topic': row['topics'].strip(),
                        'partitions': partitions,
                        'replication_factor': replication_factor,
                        'configs': configs
                    }
                    topics.append(topic_info)
            return topics
        except Exception as e:
            self.logger.error(f"Failed to read CSV file '{file}': {e}")
            return []
