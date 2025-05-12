import { check } from 'k6';
import kafka from 'k6/x/kafka';

const BOOTSTRAP_SERVERS = __ENV.BOOTSTRAP_SERVERS || 'localhost:19092,localhost:19094,localhost:19096';
const TOPIC_NAME = __ENV.TOPIC_NAME || 'load-test-topic';
const MESSAGE_SIZE_BYTES = parseInt(__ENV.MESSAGE_SIZE_BYTES || '1024');
const PRODUCER_ACKS = __ENV.PRODUCER_ACKS || 'all';
const NUM_PARTITIONS = parseInt(__ENV.NUM_PARTITIONS || '6');
const REPLICATION_FACTOR = parseInt(__ENV.REPLICATION_FACTOR || '3');
const PRODUCER_BATCH_SIZE = parseInt(__ENV.PRODUCER_BATCH_SIZE || '16384');
const PRODUCER_LINGER_MS = parseInt(__ENV.PRODUCER_LINGER_MS || '10');
const CONSUMER_GROUP = __ENV.CONSUMER_GROUP || 'load-test-consumer-group';
const CONSUMER_COUNT = parseInt(__ENV.CONSUMER_COUNT || '3');

function generateMessage() {
  return JSON.stringify({
    timestamp: Date.now(),
    sensorId: `sensor-${Math.floor(Math.random() * 100)}`,
    temperature: (20 + Math.random() * 5).toFixed(1),
    status: ['OK', 'WARN', 'ERROR', 'FATAL'][Math.floor(Math.random() * 2)]
  });
}

// producer
export async function producer() {
  const writer = new kafka.Writer({
    brokers: BOOTSTRAP_SERVERS.split(','),
    topic: TOPIC_NAME,
    acks: PRODUCER_ACKS,
    batchSize: PRODUCER_BATCH_SIZE,
    lingerMs: PRODUCER_LINGER_MS,
  });

  const payload = generateMessage();
  const key = `key-${Math.floor(Math.random() * 1000)}`;
  let success = false;

  try {
    await writer.produce({
      key,
      value: payload,
      headers: {
        source: 'xk6-kafka-load-test',
        timestamp: Date.now().toString()
      }
    });

    console.log(`Produced message with key=${key}, payload=${payload}`);
    success = true;

  } catch (error) {
    console.error(`Failed to produce message: ${error}`);
  } finally {
    writer.close();
  }

  check(success, {
    'message sent successfully': (s) => s === true,
  });
}


// consumer
export function consumer() {
  const reader = new kafka.Reader({
    brokers: BOOTSTRAP_SERVERS.split(','),
    topic: TOPIC_NAME,
    groupID: CONSUMER_GROUP,
    autoCommit: true,
  });

  const timeout = 5000;
  const startTime = Date.now();
  let messagesRead = 0;

  while (Date.now() - startTime < timeout) {
    const message = reader.consume();
    if (message) {
      messagesRead++;
      console.log(`Consumed message: key=${message.key}, value=${message.value}`);
    }
  }

  check(messagesRead, {
    'messages read successfully': (count) => count > 0,
  });

  reader.close();
}

// producer executor function
export default async function () {
  await producer();
}

// consumer executor function
export function consumeMessages() {
  consumer();
}
