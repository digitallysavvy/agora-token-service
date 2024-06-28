# CuRL scripts for testing token service using POST requests

RTC

```bash
curl -X POST http://localhost:8080/token/getNew \
-H "Content-Type: application/json" \
-d '{
  "tokenType": "rtc",
  "channel": "testChannel",
  "role": "publisher",
  "uid": "12345",
  "expire": 3600
}'

```

RTM

```bash
curl -X POST http://localhost:8080/token/getNew \
-H "Content-Type: application/json" \
-d '{
  "tokenType": "rtm",
  "uid": "12345",
  "expire": 3600
}'

```

CHAT

```bash
curl -X POST http://localhost:8080/token/getNew \
-H "Content-Type: application/json" \
-d '{
  "tokenType": "chat",
  "uid": "12345",
  "expire": 3600
}'

```
