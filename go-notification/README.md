# README

## Introduction
Create notification
- versi go apps
- versi api

Flow
- A send message, via API into Producer API
- Producer API will process the message and send to Kafka
- B request new notification to Consumer API
- Consumer API will consume message from Kafka and return the notification

Note
- seharusnya topic per user, sehingga ketika di consume langsung dikembalikan dan dihapus
- di consumer API, tetap dibutuhkan storage penampung (saat ini array) dari semnua yang diconsume. Karena API tidak dipanggil realtime, akan dipanggil ketika dibutuhkan.


## How to Run
- Start Kafka
```sh
# terminal 1 : kafka
$ cd kafka-docker
$ docker compose up -d
```

- Start Producer
```sh
# terminal 2 : producer api
cd producer
go run .
```

- Start Consumer
```sh
# terminal 3 : consumer api
cd consumer
go run .
```

- send data from API
producer
```sh
# terminal 4 : testing

# send the messsage from A
curl -X POST http://localhost:8001/api/message -d "fromUser=fauzi&toUser=adfaft&message=testing for notification"

# check the message has been received
bin/kafka-console-consumer.sh --topic mysample-topic --from-beginning --bootstrap-server localhost:9092

# retrieve the messsage by B. Jalankan berkali kali, karena goroutine sedang dalam proses manampung datanya
curl -x GET http://localhost:8002/api/message/adfaft
```

# NOTE
- 

