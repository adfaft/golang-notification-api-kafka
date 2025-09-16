# Tutorial

## Download and Start Kafka
- download kafka docker compose
### Via Docker Compse
```sh
$ curl -sSL \
https://raw.githubusercontent.com/bitnami/containers/main/bitnami/kafka/docker-compose.yml > docker-compose.yml
```
- replace KAFKA_CFG_ADVERTISED_LISTENERS from `PLAINTEXT://:9092`to `PLAINTEXT://localhost:9092`
```env
KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092
```
- Run Kafka `docker-compose up -d`

### Alternative via Docker Run
```sh
$ docker pull apache/kafka:4.1.0
$ docker run -d -p 9092:9092 --name kafka_broker apache/kafka:4.1.0
```

## Create Producer



## Create Consumer
- Initialized Gin
```sh
go mod init adfaft/go-notification/consumer
echo -e "package main \n\nimport \"fmt\" \n\nfunc main() { \n\n\tfmt.Println(\"hello\") \n\n}" > main.go
go mod tidy
go run .
```

- adfas
