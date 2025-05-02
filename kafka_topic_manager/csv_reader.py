import os
import csv

from logger import get_logger


class ReadTopicFromCSV:

    def __init__(self):
        self.logger = get_logger()

    def parse_topic_configs(self, configs_str):
        configs = {}
        for item in configs_str.split(','):
            if '=' in item:
                key, value = item.split('=', 1)
                configs[key.strip()] = value.strip()
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
                    configs_str = row.get('configs', '').strip().strip('"')
                    configs = {}
                    if configs_str:
                        configs = self.parse_topic_configs(configs_str)

                    topic_info = {
                        'topic': row['topics'].strip(),
                        'partitions': int(row['partitions']),
                        'replication_factor': int(row.get('replication_factor', default_replication_factor)),
                        'configs': configs
                    }
                    topics.append(topic_info)
        except Exception as e:
            self.logger.error(f"Failed to read File: {e}")
