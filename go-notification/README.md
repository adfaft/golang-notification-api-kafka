# README

## Introduction
Create notification
- versi go apps
- versi api


## How to Run
producer
```sh
# terminal 1 : go 
cd producer
go run .

# terminal 2 : testing
curl -X POST http://localhost:8001/api/message -d "fromUser=fauzi&toUser=adfaft&message=testing for notification"
```


