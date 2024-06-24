# Agora Backend Webservice
<p align="center">
  <img src="https://github.com/AgoraIO-Community/agora-token-service/workflows/Go/badge.svg?branch=main">
  <img src="https://github.com/AgoraIO-Community/agora-token-service/actions/workflows/dockerimage.yml/badge.svg?branch=main">
  <a href="https://github.com/AgoraIO-Community/agora-token-service/releases/latest"><img src="https://github.com/AgoraIO-Community/agora-token-service/actions/workflows/release.yml/badge.svg?release=latest"></a>
</p>

Written in Golang, using [Gin framework](https://github.com/gin-gonic/gin) to create a RESTful webservice for generating user tokens for use with the [Agora.io](https://www.agora.io) platform. 

Agora Advanced Guide: [Token Management](https://docs.agora.io/en/video-calling/develop/authentication-workflow).

## One-Click Deployments

| Railway | Render | Heroku |
|:-:|:-:|:-:|
| [![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/NKYzQA?referralCode=waRWUT) | [![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/AgoraIO-Community/agora-token-service) | [![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://www.heroku.com/deploy/?template=https://github.com/AgoraIO-Community/agora-token-service) |

## How to Run ##

Set the APP_ID, APP_CERTIFICATE and CORS_ALLOW_ORIGIN env variables.

```bash
cp .env.example .env
```

```bash
go run cmd/main.go
```

Without using `.env`, you can also set the environment variables as such:

```bash
APP_ID=app_id APP_CERTIFICATE=app_cert CORS_ALLOW_ORIGIN=allowed_origins go run cmd/main.go
```

---

The pre-compiled binaries are also available in [releases](https://github.com/AgoraIO-Community/agora-token-service/releases).

## Docker ##

#1. To build the container with app id and certificate: 

```bash
docker build -t agora-token-service --build-arg APP_ID=$APP_ID --build-arg APP_CERTIFICATE=$APP_CERTIFICATE --build-arg CORS_ALLOW_ORIGIN=$ALLOWED_ORIGINS .
```

#2. Run the container 

```bash
docker run agora-token-service
```
> Note: for testing locally
```bash
docker run -p 8080:8080 agora-token-service
```

## Makefile ##
Build and run the docker container using `make`. The `Makefile` simplifies the build and run process by reducing them to a single command. The `Makefile` uses the `.env` file to add the list of `--build-arg`'s to the `docker build` command and executes the `docker run` command after the build is completed. To avoid unnecessary rebuilds of the token server, the `Makefile` sets a `build_marker` target to watch the dockerfile, `.go` source code, and `.env` file. This enables a single command to build and run the container that only rebuilds as needed.

#1. Set the APP_ID, and APP_CERTIFICATE as env variables. Optionaly SERVER_PORT and CORS_ALLOW_ORIGIN can also be set as env variables, but they not required and there are defaults if they are not detected.
```bash
cp .env.example .env
```
#2. Run `make`
```bash
make
```
Execute the `cleanup` target to force a rebuild:
```bash
make cleanup
```

## Endpoints ##

### Ping ###
**endpoint structure**
```bash
/ping
```
response:
``` json
{"message":"pong"} 
```

### getToken ###

The `getToken` API endpoint allows you to generate tokens for different functionalities of the application. This section provides guidelines on how to use the `getToken` endpoint using HTTP POST requests.

### Endpoint URL

```
POST /getToken
```

### Request Body

The request body should contain a JSON payload with the required parameters for generating the tokens.

The following are the supported token types along with their required parameters:

1. **RTC Token:**

   To generate an RTC token for video conferencing, include the following parameters in the request body:

   ```js
   {
       "tokenType": "rtc",
       "channel": "your-channel-name",
       "role": "publisher",  // "publisher" or "subscriber"
       "uid": "your-uid",
       "expire": 3600 // optional: expiration time in seconds (default: 3600)
   }
   ```

2. **RTM Token:**

   To generate an RTM token for Real-Time Messaging, include the following parameters in the request body:

   ```js
   {
       "tokenType": "rtm",
       "uid": "your-uid",
       "channel": "test", // optional: passing channel gives streamchannel. wildcard "*" is an option.
       "expire": 3600 // optional: expiration time in seconds (default: 3600)
   }
   ```

3. **Chat Token:**

   To generate a chat token, include the following parameters in the request body:

   ```js
   {
       "tokenType": "chat",
       "uid": "your-uid", // optional: for generating a user-specific chat token
       "expire": 3600 // optional: expiration time in seconds (default: 3600)
   }
   ```

### Response

Upon successful generation of the token, the API will respond with an HTTP status code of `200 OK`, and the response body will contain the token in a JSON key `"token"`.

If there is an error during token generation or if the request parameters are invalid, the API will respond with an appropriate HTTP status code and an error message in the response body.

### Sample Usage

Here's an example of how to use the `getToken` API endpoint with a POST request using cURL:

#### Request:

```bash
curl -X POST -H "Content-Type: application/json" -d '{
    "tokenType": "rtc",
    "channel": "my-video-channel",
    "role": "publisher",
    "uid": "user123",
    "expire": 3600
}' "https://your-api-domain.com/getToken"
```

#### Reponse:

```json
{
  "token": "007hbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhZG1pbiIsInN1YiI6InVzZXIxMjMiLCJpYXQiOjE2MzEwNTU4NzIsImV4cCI6MTYzMTA1OTQ3Mn0.3eJ-RGwIl2ANFbdv4SeHtWzGiv6PpC3i0UqXlHfsqEw"
}
```

---

## Deprecated Methods
The following methods are deprecated but still operational. While they continue to work for backward compatibility, it is advised to refrain from using them in new implementations due to potential future removal or replacement with more efficient alternatives.


### RTC Token ###
The `rtc` token endpoint requires a `tokenType` (uid || userAccount), `channelName`, and the user's `uid` (type varies based on `tokenType`). 
`expiry(optional)` Pass an integer to represent the token lifetime in seconds.

**endpoint structure** 
```
/rtc/:channelName/:role/:tokenType/:rtcuid/?expiry=3600
```

response:
``` json
{"rtcToken":"007rtc-token-djfkaljdla"} 
```

### RTM Token ###

The `rtm` token endpoint requires the user's `rtmuid`. 
`expiry(optional)` Pass an integer to represent the privelege lifetime in seconds.
**endpoint structure** 
```
/rtm/:rtmuid/?expiry=3600
```

response:
``` json
{"rtmToken":"007rtm-token-djfkaljdla"} 
```

### RTM + RTC Tokens ###
The `rte` token endpoint generates both the `rtc` and `rtm` tokens with a single request. This endpoint requires a `tokenType` (uid || userAccount), `channelName`, the user's `rtcuid` (type varies `String/Int` based on `tokenType`) and `rtmuid` which is a `String`. Omitting `rtmuid` will assume it's the same as `rtcuid`.
`expiry(optional)` Pass an integer to represent the token lifetime in seconds.

**endpoint structure** 
```
/rte/:channelName/:role/:tokenType/:rtcuid/:rtmuid/?expiry=3600
```

response:
``` json
{
  "rtcToken":"007rtc-token-djfkaljdla",
  "rtmToken":"007rtm-token-djfkaljdla" 
} 
```

### Chat Tokens ###

#### endpoint structure ####

app privileges:
```
chat/app/?expiry=3600
```

user privileges:
```
/chat/account/:chatid/?expiry=3600
```

`expiry` is an optional parameter for both.

response:
``` json
{
  "chatToken":"007chat-token-djfkaljdla"
} 
```

## Contributions

Contributions are welcome, please test any changes to the Go code with the following command:

```sh
APP_ID=<YOUR_APP_ID> APP_CERTIFICATE=<YOUR_APP_CERT> go test -cover github.com/AgoraIO-Community/agora-token-service/service
```
