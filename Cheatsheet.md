# Kafka Cheatsheet

## Run Kafka
- Via Binaries
```sh
$ KAFKA_CLUSTER_ID="$(bin/kafka-storage.sh random-uuid)"
# check id berhasil digenerate
$ echo $KAFKA_CLUSTER_ID

# format log directories
$ bin/kafka-storage.sh format --standalone -t $KAFKA_CLUSTER_ID -c config/server.properties

# start the kafka
$ bin/kafka-server-start.sh config/server.properties
```

- Via Docker Image based on JVM
```sh
$ docker pull apache/kafka:4.1.0
$ docker run -p 9092:9092 apache/kafka:4.1.0
```
- Via Docker Image based on GraalVM
```sh
$ docker pull apache/kafka-native:4.1.0
$ docker run -p 9092:9092 apache/kafka-native:4.1.0
```

# Command

## Topic
- topic create , `bin/kafka-topics.sh --create --topic mysample-topic --bootstrap-server localhost:9092`
- topic describe , `bin/kafka-topics.sh --describe --topic mysample-topic --bootstrap-server localhost:9092`

## Producer
- crete event by line, `bin/kafka-console-producer.sh --topic mysample-topic --bootstrap-server localhost:9092`

## Consumer
- consume event, `bin/kafka-console-consumer.sh --topic mysample-topic --bootstrap-server localhost:9092`
- consume event from the beginning, `bin/kafka-console-consumer.sh --topic mysample-topic --from-beginning --bootstrap-server localhost:9092`