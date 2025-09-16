# Quickstart Tutorial

## Installation dan Quick Start
ref: https://kafka.apache.org/quickstart
- Get Kafka : https://www.apache.org/dyn/closer.cgi?path=/kafka/4.1.0/kafka_2.13-4.1.0.tgz and extract

NOTES: untuk **path directory project jangan ada space**, kalau adea akan error
`Classpath is empty. Please build the project first e.g. by running './gradlew jar -PscalaVersion=2.13.16'`

```sh
$ tar -xzf kafka_2.13-4.1.0.tgz
$ cd kafka_2.13-4.1.0

# test binary berhasil
$ bin/kafka-storage.sh random-uuid
```

- Start Environemnt
```sh
# generate cluster uuid
$ KAFKA_CLUSTER_ID="$(bin/kafka-storage.sh random-uuid)"
# check id berhasil digenerate
$ echo $KAFKA_CLUSTER_ID

# format log directories
$ bin/kafka-storage.sh format --standalone -t $KAFKA_CLUSTER_ID -c config/server.properties

# start the kafka
$ bin/kafka-server-start.sh config/server.properties
# tunggu sampai "server started" message tampil 
```

Alternatively, using docker image based on JVM Based
```sh
$ docker pull apache/kafka:4.1.0
$ docker run -p 9092:9092 apache/kafka:4.1.0
```

Alternatively, using docker image based on GraalVM Based Natifve Apache Kakfka 
```sh
$ docker pull apache/kafka-native:4.1.0
$ docker run -p 9092:9092 apache/kafka-native:4.1.0
```

- Create Topic untuk menyimpan event
```sh
# di terminal baru terpisah (AKA terminal 2)
$ bin/kafka-topics.sh --create --topic mysample-topic --bootstrap-server localhost:9092
# akan menampilkan topic berhasil dibuat

# describe / jelaskan isi dari specific topic seperti ID, partiition, replica, etc
$ bin/kafka-topics.sh --describe --topic mysample-topic --bootstrap-server localhost:9092

```

- Create an event using the topic
```sh
# create an event, per line
$ bin/kafka-console-producer.sh --topic mysample-topic --bootstrap-server localhost:9092
```

- Consume an event using the topic
```sh
# on another terminal (AKA terminal 3, recommend as split terminal)
$ bin/kafka-console-consumer.sh --topic mysample-topic --from-beginning --bootstrap-server localhost:9092
# jika producer di terminal 2 ditambah informasinya, maka consumer akan otomatis menampilkan data terbaru juga
```

- Finish, now teardown the apps dengan CTRL+C untuk producer dan semua consumer serta Kafka Broker sendiri. 
Jangan lupa untuk delete temporary folder `rm -rf /tmp/kafka-logs /tmp/kraft-combined-logs`

## Optional
### Kafka Connect
ref: https://kafka.apache.org/quickstart#quickstart_kafkaconnect
Kafka Connect digunakan untuk membaca data dari external untuk dikirimnkan ke Kafka 

### Kafka Streams
ref: https://kafka.apache.org/41/documentation/streams/quickstart
Kafka Steams digunakan untuk membaca/menulis kafka stream. Contoh, dengan library ini via Java atau Scala, dia membaca sebagai consumer, kemudian kita olah data  dan kirimkan lagi ke topic berbeda sebagai producer

# Ask
- How to change kafka default port from 9092 ?
- How to consume not from beginning ?
  - Jika di consumer tidak menggunakan flag `--from-beginning` maka akan dibaca dari terakhir topic dibaca
- How multiple consumer read the data if the consumer start the application from diferent time?
  - Maka consumer tetap dibaca dari terakhir topic itu di consume, tidak peduli dari multiple apps